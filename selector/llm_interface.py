from abc import ABC, abstractmethod
from typing import List
from models import Tool

class LLMServiceInterface(ABC):
    
    @abstractmethod
    async def select_best_tool(self, user_prompt: str, candidate_tools: List[Tool]) -> str:
        pass
    
    @abstractmethod
    async def generate_selection_message(self, user_prompt: str, selected_tool: Tool) -> str:
        pass 