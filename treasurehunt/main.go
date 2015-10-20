package main

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Zumata/grpc-challenge/treasurehunt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "treasurehunt.zumata.com:3000"
	defaultName = "test45"
)

func main() {
	go func() {
		http.ListenAndServe(":2016", nil)
	}()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := treasurehunt.NewTreasureHuntClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	ctx := context.Background()
	r, err := c.Join(ctx, &treasurehunt.JoinRequest{Name: name})
	if err != nil {
		log.Fatalf("could not join: %v", err)
	}
	md := metadata.New(map[string]string{"id": r.Id})
	ctx = metadata.NewContext(ctx, md)
	log.Printf("Joined: %#+v", r)
	log.Printf("Position: %#+v", *r.Position)

	curPos := *r.Position
	for {
		t, err := c.GetTreasure(ctx, &treasurehunt.GetTreasureRequest{})
		time.Sleep(200 * time.Millisecond)
		if err != nil {
			log.Fatalf("Failed to get treasure: %v", err)
		}
		minDist := int64(math.MaxInt64)
		var target *treasurehunt.Position
		for _, pos := range t.Position {
			dist := abs(pos.X-curPos.X) + abs(pos.Y-curPos.Y)
			if dist < minDist {
				target = pos
				minDist = dist
			}
		}
		if target == nil {
			continue
		}
		log.Printf("Going to: %#+v", *target)
		if target.X != curPos.X {
			dir := treasurehunt.MoveRequest_LEFT
			if target.X > curPos.X {
				dir = treasurehunt.MoveRequest_RIGHT
			}
			log.Printf("pos: %v, dir: %v", curPos, dir)
			m, err := c.Move(ctx, &treasurehunt.MoveRequest{Direction: dir})
			if err != nil {
				log.Printf("Move failed: %v", err)
			} else {
				log.Printf("Moved: %+#v", m)
				curPos = *m.Position
				log.Printf("Curpos: %+#v", curPos)
			}
			time.Sleep(200 * time.Millisecond)
		}

		if target.Y != curPos.Y {
			dir := treasurehunt.MoveRequest_DOWN
			if target.Y > curPos.Y {
				dir = treasurehunt.MoveRequest_UP
			}
			m, err := c.Move(ctx, &treasurehunt.MoveRequest{Direction: dir})
			time.Sleep(200 * time.Millisecond)
			if err != nil {
				log.Printf("Move failed: %v", err)
			} else {
				log.Printf("Moved: %+#v", m)
				curPos = *m.Position
				log.Printf("Curpos: %+#v", curPos)
			}
		}

	}

}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
