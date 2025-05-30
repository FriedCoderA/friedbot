package triggers

import "friedbot/pkg/models/schema"

type Scorer interface {
	score(session *schema.Session, score int) int
}
