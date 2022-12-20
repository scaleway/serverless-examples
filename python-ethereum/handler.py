#!/usr/bin/env python3.9
# -*- coding: utf-8 -*-

import sys, secrets, time, json, msgpack, argon2, base64
import string

from eth_account import Account

def handler(event, context):
    """Handle a request to the Lyra serverless funciton
    Args:
        event (dict): request params
        context (dict): function call metadata
    
    Description:
        * generates a private key from true random numbers
        * attaches 0x to make the hexadecimal look like an address
        * computes an ethereum address
        * calculates a unix timestamp
        * transforms into subbinary messagepack serialization
        * creates an argon2 hash of the serialization with salt 'scaleway'
        * stores hashsum, payload, and obfuscated salt into a json
        * theoretic concept of signing a data structure with a mutually known secret
    
    """

    _salt_ = str(base64.b64decode(b'c2NhbGV3YXk=')) # technically argon2 is not sha1 ... 
    _binarysecret_ = 0x01
    priv = secrets.token_hex(32)
    private_key = "0x" + priv
    acct = Account.from_key(private_key)

    payload = {
        'privkey'   : private_key,
        'address'   : acct.address,
        'timestamp' : time.time(),
    }

    string = msgpack.packb(payload, use_bin_type=True)
    hash = argon2.argon2_hash(string, _salt_)

    data = {
        'hashsum':    str(base64.b64encode(hash)),
        'payload':    payload,
        'transport':  str(base64.b64encode(bytes(_salt_, "utf-8"))),
    }

    return {
        "body": json.dumps(data), 
        "headers": {
            "Content-Type": ["text/json"],
        }
    }
