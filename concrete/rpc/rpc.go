package rpc

import "github.com/ethereum/go-ethereum/eth"

type ConcreteRPC struct {
	Namespace     string      // namespace under which the rpc methods of Service are exposed
	Service       interface{} // receiver instance which holds the methods
	Authenticated bool        // whether the api should only be available behind authentication.
}

type ConcreteRPCConstructor func(*eth.Ethereum) ConcreteRPC
