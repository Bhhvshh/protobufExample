import pymongo
import pika
import json
import time
from datetime import datetime

def encode_datetime(obj):
    if isinstance(obj, datetime):
        return obj.isoformat()
    raise TypeError(f"Object of type {type(obj)} is not JSON serializable")

def main():
    # MongoDB connection
    mongo_client = pymongo.MongoClient('mongodb://localhost:27017/')
    db = mongo_client['grpc_messages']
    collection = db['greetings']

    # RabbitMQ connection
    rabbitmq_connection = pika.BlockingConnection(
        pika.ConnectionParameters(host='localhost')
    )
    channel = rabbitmq_connection.channel()
    channel.queue_declare(queue='greetings')

    # Get the last processed timestamp
    last_timestamp = datetime.min

    try:
        while True:
            # Find new documents
            query = {'timestamp': {'$gt': last_timestamp}}
            cursor = collection.find(query).sort('timestamp', pymongo.ASCENDING)
            
            for doc in cursor:
                # Update last processed timestamp
                last_timestamp = doc['timestamp']
                
                # Remove MongoDB _id before sending
                doc_without_id = {k: v for k, v in doc.items() if k != '_id'}
                
                # Send to RabbitMQ
                message = json.dumps(doc_without_id, default=encode_datetime)
                channel.basic_publish(
                    exchange='',
                    routing_key='greetings',
                    body=message
                )
                print(f"Sent message to RabbitMQ: {message}")
            
            # Wait before next check
            time.sleep(1)
            
    except KeyboardInterrupt:
        print("Stopping CDC service...")
    finally:
        rabbitmq_connection.close()
        mongo_client.close()

if __name__ == '__main__':
    main()