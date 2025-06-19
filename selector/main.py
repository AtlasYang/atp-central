from fastapi import FastAPI
from api import router
import uvicorn
from huggingface_hub import login
import os

login(token=os.environ["HUGGINGFACE_TOKEN"])

app = FastAPI(title="Tool Selector Service", version="1.0.0")

# Include routes
app.include_router(router, prefix="/api/v1")

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=int(os.environ["PORT"]),
        reload=True
    )