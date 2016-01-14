package algo

type Hyperparams struct {
	TopicPrior    []float64
	TopicPriorSum float64
	WordPrior     float64
	WordPriorSum  float64
	VocabSize     int32
}

func (h *Hyperparams) Topics() int {
	return len(h.TopicPrior)
}

func (h *Hyperparams) OptimizeTopicPrior(docLenCount []int32, topicDocCount [][]int32, shape, scale float64, iterations int) {
	for iter := 0; iter < iterations; iter++ {
		denom := 0.0
		diffDigamma := 0.0
		for j := 1; j < len(docLenCount); j++ {
			diffDigamma += 1.0 / (float64(j-1) + h.TopicPriorSum)
			denom += float64(docLenCount[j]) * diffDigamma
		}
		denom -= 1.0 / scale

		h.TopicPriorSum = 0
		for k := 0; k < len(topicDocCount); k++ {
			numer := 0.0
			diffDigamma = 0.0
			for j := 1; j < len(topicDocCount); j++ {
				diffDigamma += 1.0 / (float64(j-1) + h.TopicPrior[k])
				numer += float64(topicDocCount[k][j]) * diffDigamma
			}
			h.TopicPrior[k] = (h.TopicPrior[k]*numer + shape) / denom
			h.TopicPriorSum += h.TopicPrior[k]
		}
	}
}

func (h *Hyperparams) OptimizeWordPrior(topicLenCount, wordTopicCount []int32, iterations int) {
	for iter := 0; iter < iterations; iter++ {
		numer := 0.0
		diffDigamma := 0.0
		for j := 1; j < len(wordTopicCount); j++ {
			diffDigamma += 1.0 / (float64(j-1) + h.WordPrior)
			numer += diffDigamma * float64(wordTopicCount[j])
		}

		denom := 0.0
		diffDigamma = 0.0
		for j := 1; j < len(topicLenCount); j++ {
			diffDigamma += 1.0 / (float64(j-1) + h.WordPriorSum)
			denom += diffDigamma * float64(topicLenCount[j])
		}
		h.WordPriorSum = h.WordPrior * numer / denom
		h.WordPrior = h.WordPriorSum / float64(h.VocabSize)
	}
}
