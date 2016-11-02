/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"log"
	"os"
	"net/http"
	"io/ioutil"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/stevenmatthewt/GRPC_REST_Test/helloworld"
)

const (
	RESTaddress = "http://localhost:8080"
	GRPCaddress = "localhost:50051"
	defaultName = "world"
	grpcPort    = 50051
)

func main() {
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	contactGRPC(name)
	contactREST(name)
}

func contactGRPC(name string) {
	log.Print("===== GRPC =====")
	// Set up a connection to the server.
	conn, err := grpc.Dial(GRPCaddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	hello, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", hello.Message)

	goodbye, err := c.SayGoodbye(context.Background(), &pb.GoodbyeRequest{Name: name})
	if err != nil {
		log.Fatalf("could not say goodbye: %v", err)
	}
	log.Printf("Response: %s", goodbye.Message)
}

func contactREST(name string) {
	log.Print("===== REST =====")

	resp, err := http.Get(RESTaddress + "/hello/Steven")
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("could not parse response: %v", err)
		return
	}
	log.Printf("Response: %s", string(bytes))

	resp2, err := http.Get(RESTaddress + "/goodbye/Steven")
	if err != nil {
		log.Fatalf("could not say goodbye: %v", err)
		return
	}
	defer resp2.Body.Close()
	bytes, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalf("could not parse response: %v", err)
		return
	}
	log.Printf("Response: %s", string(bytes))
}
