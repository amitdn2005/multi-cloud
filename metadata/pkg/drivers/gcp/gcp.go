// Copyright 2023 The SODA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcp

import (
	"context"

	backendpb "github.com/opensds/multi-cloud/backend/proto"
	pb "github.com/opensds/multi-cloud/metadata/proto"
	log "github.com/sirupsen/logrus"
	"github.com/webrtcn/s3client"
)

type GcpAdapter struct {
	backend *backendpb.BackendDetail
	session *s3client.Client
}

func BucketList(sess *session.Session) ([]*model.MetaBucket, error) {
	svc := s3.New(sess)

	output, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Errorf("unable to list buckets. failed with error: %v", err)
		return nil, err
	}
	numBuckets := len(output.Buckets)
	log.Info("Amit: GCP ===========")
	log.Info(output.Buckets)
	bucketArray := make([]*model.MetaBucket, numBuckets)
	wg := sync.WaitGroup{}
	for i, bucket := range output.Buckets {
		wg.Add(1)
		// go GetBucketMeta(i, bucket, sess, bucketArray, &wg)
	}
	wg.Wait()
	return bucketArray, err
}

func (ad *GcpAdapter) SyncMetadata(ctx context.Context, in *pb.SyncMetadataRequest) error {
	buckArr, err := BucketList(ad.Session)
	if err != nil {
		log.Errorf("metadata collection for backend id: %v failed with error: %v", ad.Backend.Id, err)
		return err
	}

	metaBackend := model.MetaBackend{}
	metaBackend.Id = bson.ObjectId(ad.Backend.Id)
	metaBackend.BackendName = ad.Backend.Name
	metaBackend.Type = ad.Backend.Type
	metaBackend.Region = ad.Backend.Region
	metaBackend.Buckets = buckArr
	newContext := context.TODO()
	err = db.DbAdapter.CreateMetadata(newContext, metaBackend)

	return err
}
func (ad *GcpAdapter) DownloadObject() {
	log.Info("Implement me (gcp) driver")
}
