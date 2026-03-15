import os
import json
import redis

class CacheRedisConfig:
    def __init__(self, host='localhost', port=6379, db=0):
        self.host = host
        self.port = port
        self.db = db
        self.client = redis.Redis(host=host, port=port, db=db)

    def get_config(self, key):
        value = self.client.get(key)
        if value is None:
            return None
        return json.loads(value.decode('utf-8'))

    def set_config(self, key, value):
        self.client.set(key, json.dumps(value).encode('utf-8'))

    def delete_config(self, key):
        self.client.delete(key)

def main():
    config = CacheRedisConfig()
    config.set_config('database', {'host': 'localhost', 'port': 5432, 'username': 'user', 'password': 'password'})
    print(config.get_config('database'))
    config.delete_config('database')
    print(config.get_config('database'))

if __name__ == '__main__':
    main()