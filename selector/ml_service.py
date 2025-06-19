import torch
from typing import List
from transformers import AutoTokenizer, AutoModelForCausalLM, BitsAndBytesConfig
from peft import PeftModel
from models import Tool
from config import settings
from llm_interface import LLMServiceInterface

class LlamaModelService(LLMServiceInterface):
    def __init__(self):
        self.model_name = settings.model_name
        self.adapter_path = settings.adapter_path
        self.tokenizer = None
        self.model = None
        self._model_loaded = False

    def _load_model(self):
        """Load model and tokenizer"""
        if self._model_loaded:
            return
            
        print("Loading model...")
        
        self.tokenizer = AutoTokenizer.from_pretrained(
            self.model_name,
            trust_remote_code=True
        )
        
        if torch.cuda.is_available():
            bnb_config = BitsAndBytesConfig(
                load_in_4bit=True,
                bnb_4bit_use_double_quant=True,
                bnb_4bit_quant_type="nf4",
                bnb_4bit_compute_dtype=torch.bfloat16
            )
            
            base_model = AutoModelForCausalLM.from_pretrained(
                self.model_name,
                quantization_config=bnb_config,
                device_map="auto",
                trust_remote_code=True,
                torch_dtype=torch.bfloat16
            )
        else:
            base_model = AutoModelForCausalLM.from_pretrained(
                self.model_name,
                trust_remote_code=True,
                torch_dtype=torch.float32
            )
        
        self.model = PeftModel.from_pretrained(base_model, self.adapter_path)
        self._model_loaded = True
        
        print("Model loaded successfully")

    async def select_best_tool(self, user_prompt: str, candidate_tools: List[Tool]) -> str:
        """Select the most appropriate tool based on user prompt"""
        self._load_model()
        
        if not candidate_tools:
            raise ValueError("No candidate tools provided")
        
        inference_prompt = self._create_inference_prompt(user_prompt, candidate_tools)
        inputs = self.tokenizer(inference_prompt, return_tensors="pt")
        
        if torch.cuda.is_available():
            inputs = inputs.to("cuda")
        
        with torch.no_grad():
            outputs = self.model.generate(
                **inputs,
                max_new_tokens=10,
                eos_token_id=self.tokenizer.eos_token_id,
                pad_token_id=self.tokenizer.eos_token_id,
                do_sample=False
            )
        
        response = self.tokenizer.decode(
            outputs[0][inputs['input_ids'].shape[1]:], 
            skip_special_tokens=True
        ).strip()
        
        return response

    async def generate_selection_message(self, user_prompt: str, selected_tool: Tool) -> str:
        """Generate explanation message for why the tool was selected"""
        self._load_model()
        
        message_prompt = self._create_message_prompt(user_prompt, selected_tool)
        inputs = self.tokenizer(message_prompt, return_tensors="pt")
        
        if torch.cuda.is_available():
            inputs = inputs.to("cuda")
        
        with torch.no_grad():
            outputs = self.model.generate(
                **inputs,
                max_new_tokens=150,  # 더 긴 응답을 위해 토큰 수 증가
                eos_token_id=self.tokenizer.eos_token_id,
                pad_token_id=self.tokenizer.eos_token_id,
                do_sample=True,
                temperature=0.7,
                top_p=0.9
            )
        
        response = self.tokenizer.decode(
            outputs[0][inputs['input_ids'].shape[1]:], 
            skip_special_tokens=True
        ).strip()
        
        return response

    def _create_inference_prompt(self, user_prompt: str, candidate_tools: List[Tool]) -> str:
        """Create inference prompt for tool selection"""
        tool_info = []
        for tool in candidate_tools:
            tool_info.append(f"- {tool.name}: {tool.description or 'No description available'}")
        
        context_str = "[Available Tools]\n" + "\n".join(tool_info)
        instruction_str = f"Based on the provided tool descriptions and user question, select the most appropriate tool.\n\n[User Question]\n{user_prompt}"
        
        messages = [
            {"role": "user", "content": f"{instruction_str}\n\n{context_str}"}
        ]
        
        return self.tokenizer.apply_chat_template(
            messages, 
            tokenize=False, 
            add_generation_prompt=True
        )

    def _create_message_prompt(self, user_prompt: str, selected_tool: Tool) -> str:
        """Create prompt for explaining why the tool was selected"""
        instruction_str = f"""You are an AI assistant that explains tool selections in conversational manner.

A user asked: "{user_prompt}"

The tool "{selected_tool.name}" is selected for this task.

Tool description: {selected_tool.description or 'No description available'}

Please explain in a natural, conversational way why this tool is appropriate for the user's request. Keep it brief but informative.
Use language of the user's request.
"""
        
        messages = [
            {"role": "user", "content": instruction_str}
        ]
        
        return self.tokenizer.apply_chat_template(
            messages, 
            tokenize=False, 
            add_generation_prompt=True
        )

model_service = LlamaModelService() 