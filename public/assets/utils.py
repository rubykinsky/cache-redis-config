import os
import json
import logging
from typing import Dict, Any

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def load_config(file_path: str) -> Dict[str, Any]:
    try:
        with open(file_path, 'r') as file:
            return json.load(file)
    except FileNotFoundError:
        logger.error(f"Config file not found: {file_path}")
        return {}
    except json.JSONDecodeError as e:
        logger.error(f"Failed to parse config file: {e}")
        return {}

def save_config(file_path: str, config: Dict[str, Any]) -> None:
    try:
        with open(file_path, 'w') as file:
            json.dump(config, file, indent=4)
    except Exception as e:
        logger.error(f"Failed to save config file: {e}")

def get_redis_config(config: Dict[str, Any]) -> Dict[str, Any]:
    return config.get('redis', {})

def update_redis_config(config: Dict[str, Any], redis_config: Dict[str, Any]) -> Dict[str, Any]:
    config['redis'] = redis_config
    return config

def get_env_var(var_name: str) -> str:
    return os.getenv(var_name)

def main():
    config_file = get_env_var('CONFIG_FILE')
    if config_file:
        config = load_config(config_file)
        redis_config = get_redis_config(config)
        logger.info(f"Loaded Redis config: {redis_config}")
    else:
        logger.error("CONFIG_FILE environment variable not set")

if __name__ == '__main__':
    main()