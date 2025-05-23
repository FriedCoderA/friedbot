package triggers

import "friedbot/pkg/models/schema"

type Trigger interface {
	score(session *schema.Session, score int) int
}
