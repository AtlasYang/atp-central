from dotenv import load_dotenv
import os

load_dotenv()

class Settings:
    def __init__(self):
        self.database_url = os.getenv("MAIN_DB_CONNECTION")

        self.llm_provider = os.getenv("LLM_PROVIDER", "standalone").lower()  # "standalone" or "openai"

        self.model_name = "meta-llama/Llama-3.1-8B-Instruct"
        self.adapter_path = "./sft-final-adapter"

        self.openai_api_key = os.getenv("OPENAI_API_KEY")
        self.openai_model = os.getenv("OPENAI_MODEL", "gpt-4o-mini")

    @property
    def db_connection_string(self) -> str:
        if not self.database_url:
            raise ValueError("MAIN_DB_CONNECTION environment variable is not set")
        return self.database_url

settings = Settings() 