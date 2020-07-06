from __future__ import print_function
import logging

import json
import grpc

import user_pb2
import user_pb2_grpc

def connect():
    channel = grpc.insecure_channel('localhost:9000')
    stub = user_pb2_grpc.UserServiceStub(channel)
    return stub

def saveUsers(stub):
    for i in range(1, 11):
        data = openFile(i)
        for i in range(len(data)):
            saveUser(stub, data[i])

def saveUser(stub, user):
    response = stub.Save(user_pb2.User(
        id=user['id'],
        firstName=user['first_name'],
        lastName=user['last_name'],
        email=user['email'],
        gender=user['gender'],
        ipAddress=user['ip_address'],
        userName=user['user_name'],
        agent=user['agent'],
        country=user['country']
    ))

def openFile(i):
    fileName = str(i) + ".json"
    with open('client/users/' + fileName) as json_file:
        data = json.load(json_file)
        return data

if __name__ == '__main__':
    stub = connect()
    saveUsers(stub)
