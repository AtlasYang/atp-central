from openai import AsyncOpenAI
from typing import List
from models import Tool
from config import settings
from llm_interface import LLMServiceInterface

class OpenAIModelService(LLMServiceInterface):
    def __init__(self):
        if not settings.openai_api_key:
            raise ValueError("OPENAI_API_KEY environment variable is not set")
        
        self.client = AsyncOpenAI(api_key=settings.openai_api_key)
        self.model = settings.openai_model
    
    async def select_best_tool(self, user_prompt: str, candidate_tools: List[Tool]) -> str:
        """Select the most appropriate tool based on user prompt using OpenAI GPT"""
        if not candidate_tools:
            raise ValueError("No candidate tools provided")
        
        tool_info = []
        for tool in candidate_tools:
            tool_info.append(f"- {tool.name}: {tool.description or 'No description available'}")
        
        context_str = "[Available Tools]\n" + "\n".join(tool_info)
        instruction_str = f"Based on the provided tool descriptions and user question, select the most appropriate tool. Reply with only the tool name.\n\n[User Question]\n{user_prompt}"
        
        try:
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[
                    {"role": "system", "content": "You are a helpful assistant that selects the most appropriate tool for a given task. Respond only with the tool name."},
                    {"role": "user", "content": f"{instruction_str}\n\n{context_str}"}
                ],
                max_tokens=50,
                temperature=0.1
            )
            
            return response.choices[0].message.content.strip()
        except Exception as e:
            raise ValueError(f"Failed to select tool using OpenAI: {str(e)}")
    
    async def generate_selection_message(self, user_prompt: str, selected_tool: Tool) -> str:
        """Generate explanation message for why the tool was selected using OpenAI GPT"""
        instruction_str = f"""You are an AI assistant that explains tool selections in conversational manner.

A user asked: "{user_prompt}"

The tool "{selected_tool.name}" is selected for this task.

Tool description: {selected_tool.description or 'No description available'}

Please explain in a natural, conversational way why this tool is appropriate for the user's request. Keep it brief but informative.
Use language of the user's request.
"""
        
        try:
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[
                    {"role": "system", "content": "You are a helpful assistant that explains tool selections in a friendly, conversational manner."},
                    {"role": "user", "content": instruction_str}
                ],
                max_tokens=200,
                temperature=0.7
            )
            
            return response.choices[0].message.content.strip()
        except Exception as e:
            raise ValueError(f"Failed to generate message using OpenAI: {str(e)}")

# Factory function to create the appropriate service
def create_llm_service() -> LLMServiceInterface:
    """Factory function to create the appropriate LLM service based on configuration"""
    if settings.llm_provider == "openai":
        return OpenAIModelService()
    elif settings.llm_provider == "standalone":
        from ml_service import LlamaModelService
        return LlamaModelService()
    else:
        raise ValueError(f"Unknown LLM provider: {settings.llm_provider}")

# Create the service instance based on configuration
model_service = create_llm_service() 