package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_ParseMilvuseContent(t *testing.T) {
	src := `"{\"content\":\"\n MessagingThe core messaging capabilities of NATS includingpub-sub,request-reply, andqueue groups.Core Publish-SubscribeRequest-ReplyJSON for Message PayloadsProtobuf for Message PayloadsConcurrent Message ProcessingIterating Over Multiple SubscriptionsJetStreamAnintegrated subsystemproviding distributed persistence.Limits-based StreamInterest-based StreamWork-queue StreamPull ConsumersPull Consumer - Applying LimitsPush Consumers (legacy)Queue Push Consumers (legacy)Multi-Stream Consumption (legacy)Migration to new JetStream APIConsumer - Fetch MessagesSubject-Mapped PartitionsAuthentication and AuthorizationTopics related toauthenticationandauthorization.Programmatic NKeys and JWTsConfiguring the System AccountPrivate InboxAuth Callout - CentralizedAuth Callout - DecentralizedPrivate Inbox using JWTTopologiesExamples showcasing various deployment topologies NATS supports includingclusters,superclusters, andleaf nodes.Simple LeafnodeLeafnode with JWT AuthSuperclusterSupercluster with JetStreamSupercluster ArbiterUse CasesCross-functional examples satisfying a use case.Regional and Cross Region Streams (Supercluster)Regional and Cross Region Streams (Cluster)Key-ValueA layer on top of JetStream for utilizing a stream as\"}"`

	contentS := ""
	err := json.Unmarshal([]byte(src), &contentS)
	assert.Nil(t, err)
	contentS = strings.ReplaceAll(contentS, "\n", "")
	content := make(map[string]string)
	err = json.Unmarshal([]byte(contentS), &content)
	t.Log(content["content"])
	t.Log(contentS)

}
