package algo

// For algorithms used please refer to "Efficient Methods for Topic
// Model Inference on Streaming Document Collections" (KDD 2009).
/*
  TrainSamplerBase: model_;
  HyperparamsOptimTrainSampler: docLenCount, topicDocCount, topicLenCount, wordTopicCount
  SparseLDATrainSampler: smoothingOnlyBucket, docTopicBucket, wordTopicBucket, cachedCoeff
*/
type Sampler struct {
	model *Model
}
