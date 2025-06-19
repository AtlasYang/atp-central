from dataclasses import dataclass
from typing import Optional
from datetime import datetime
from pydantic import BaseModel

@dataclass
class Tool:
    id: int
    uuid: str
    name: str
    version: str
    description: Optional[str]
    engine_interface: str
    provider_interface: str
    created_at: datetime
    updated_at: datetime

class SelectRequest(BaseModel):
    user_prompt: str
    
class SelectResponse(BaseModel):
    tool_id: int
    message: str