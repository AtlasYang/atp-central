from typing import List
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np

from models import Tool, SelectRequest, SelectResponse
from database import tool_repository
from openai_service import model_service

class ToolSelectorService:
    def __init__(self):
        self.tool_repository = tool_repository
        self.model_service = model_service
        self.embedding_model = SentenceTransformer('all-MiniLM-L6-v2')

    async def select_tool(self, request: SelectRequest) -> SelectResponse:
        """Select the most appropriate tool for user prompt"""
        all_tools = await self.tool_repository.get_all_tools()
        
        if not all_tools:
            raise ValueError("No tools available")
        
        candidate_tools = self._select_candidate_tools(request.user_prompt, all_tools)
        
        selected_tool_name = await self.model_service.select_best_tool(
            request.user_prompt, candidate_tools
        )
        
        selected_tool = await self.tool_repository.get_tool_by_name(selected_tool_name)
        
        explanation_message = await self.model_service.generate_selection_message(
            request.user_prompt, selected_tool
        )
        
        return SelectResponse(tool_id=selected_tool.id, message=explanation_message)

    def _select_candidate_tools(self, user_prompt: str, all_tools: List[Tool], top_k: int = 5) -> List[Tool]:
        """Select top K similar tools using SentenceTransformer embedding similarity"""
        if len(all_tools) <= top_k:
            return all_tools
            
        tool_descriptions = [tool.description or "" for tool in all_tools]
        tool_embeddings = self.embedding_model.encode(tool_descriptions)
        user_embedding = self.embedding_model.encode([user_prompt])[0]
        
        similarities = cosine_similarity([user_embedding], tool_embeddings)[0]
        
        top_indices = np.argsort(similarities)[-top_k:][::-1]
        return [all_tools[i] for i in top_indices]

tool_selector_service = ToolSelectorService() 