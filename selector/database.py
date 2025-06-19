import psycopg2
from typing import List
from psycopg2.extras import RealDictCursor
from models import Tool
from config import settings

class ToolRepository:
    def __init__(self):
        self.connection_string = settings.db_connection_string

    def _get_connection(self):
        return psycopg2.connect(self.connection_string)

    async def get_all_tools(self) -> List[Tool]:
        """Retrieve all tools from database"""
        with self._get_connection() as conn:
            with conn.cursor(cursor_factory=RealDictCursor) as cursor:
                cursor.execute("""
                    SELECT id, uuid, name, version, description, 
                           engine_interface, provider_interface, 
                           created_at, updated_at
                    FROM tools 
                    ORDER BY created_at DESC
                """)
                rows = cursor.fetchall()
                
                return [
                    Tool(
                        id=row['id'],
                        uuid=row['uuid'],
                        name=row['name'],
                        version=row['version'],
                        description=row['description'],
                        engine_interface=row['engine_interface'],
                        provider_interface=row['provider_interface'],
                        created_at=row['created_at'],
                        updated_at=row['updated_at']
                    )
                    for row in rows
                ]

    async def get_tool_by_name(self, name: str) -> Tool:
        """Retrieve tool by name"""
        with self._get_connection() as conn:
            with conn.cursor(cursor_factory=RealDictCursor) as cursor:
                cursor.execute("""
                    SELECT id, uuid, name, version, description, 
                           engine_interface, provider_interface, 
                           created_at, updated_at
                    FROM tools 
                    WHERE name = %s
                    LIMIT 1
                """, (name,))
                row = cursor.fetchone()
                
                if not row:
                    raise ValueError(f"Tool with name '{name}' not found")
                
                return Tool(
                    id=row['id'],
                    uuid=row['uuid'],
                    name=row['name'],
                    version=row['version'],
                    description=row['description'],
                    engine_interface=row['engine_interface'],
                    provider_interface=row['provider_interface'],
                    created_at=row['created_at'],
                    updated_at=row['updated_at']
                )

tool_repository = ToolRepository() 