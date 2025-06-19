from fastapi import APIRouter, HTTPException
from models import SelectRequest, SelectResponse
from services import tool_selector_service

router = APIRouter()

@router.post("/select", response_model=SelectResponse)
async def select_tool(request: SelectRequest):
    """Select the most appropriate tool for the given user prompt"""
    try:
        response = await tool_selector_service.select_tool(request)
        return response
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}") 