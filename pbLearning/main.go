package main

import (
	"fmt"
	"log"

	pb "pbLearning/proto"

	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("hello")

	bruce := &pb.Person{
		Name: "bruce",
		Age:  23,
		Socialfollowers: &pb.Socialfollowers{
			Youtube: 100,
			Twitter: 200,
		},
	}

	data, err := proto.Marshal(bruce)
	if err != nil {
		log.Fatal("Marshalling error:", err)
	}
	fmt.Println(data)

	newBruce := &pb.Person{}
	err = proto.Unmarshal(data, newBruce)
	if err != nil {
		log.Fatal("Unmarshalling error:", err)
	}
	fmt.Println(newBruce.GetAge())
	fmt.Println(newBruce.GetName())
	fmt.Println(newBruce.Socialfollowers.GetTwitter())
	fmt.Println(newBruce.Socialfollowers.GetYoutube())
}
