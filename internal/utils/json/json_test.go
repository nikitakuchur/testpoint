package json_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/utils/json"
	"testing"
)

func TestReformatJson(t *testing.T) {
	input := `{"id":"DA8GpfGvqqNf3nUQMfcCwA","words":[{"id":377534,"text":"spelling","pronunciationTracks":[{"id":44,"variety":"UK","filepath":"/pronunciations/en-GB-Standard-A-spelling.mp3"},{"id":43,"variety":"US","filepath":"/pronunciations/en-US-Standard-G-spelling.mp3"},{"id":41,"variety":"AUS","filepath":"/pronunciations/en-AUS-Standard-B-spelling.mp3"}]},{"id":304903,"text":"popular","pronunciationTracks":[{"id":32,"variety":"US","filepath":"/pronunciations/en-US-Standard-A-popular.mp3"},{"id":30,"variety":"US","filepath":"/pronunciations/en-AUS-Standard-B-popular.mp3"},{"id":37,"variety":"AUS","filepath":"/pronunciations/en-GB-Standard-G-popular.mp3"}]},{"id":413482,"text":"train","pronunciationTracks":[{"id":55,"variety":"AUS","filepath":"/pronunciations/en-AUS-Standard-A-train.mp3"},{"id":58,"variety":"US","filepath":"/pronunciations/en-US-Standard-G-train.mp3"},{"id":52,"variety":"UK","filepath":"/pronunciations/en-GB-Standard-B-train.mp3"}]}],"timestamp":"2024-06-05T22:51:09.464637Z"}`

	actual := json.ReformatJson(input, false, []string{})
	expected := `{
  "id": "DA8GpfGvqqNf3nUQMfcCwA",
  "timestamp": "2024-06-05T22:51:09.464637Z",
  "words": [
    {
      "id": 377534,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-GB-Standard-A-spelling.mp3",
          "id": 44,
          "variety": "UK"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-G-spelling.mp3",
          "id": 43,
          "variety": "US"
        },
        {
          "filepath": "/pronunciations/en-AUS-Standard-B-spelling.mp3",
          "id": 41,
          "variety": "AUS"
        }
      ],
      "text": "spelling"
    },
    {
      "id": 304903,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-US-Standard-A-popular.mp3",
          "id": 32,
          "variety": "US"
        },
        {
          "filepath": "/pronunciations/en-AUS-Standard-B-popular.mp3",
          "id": 30,
          "variety": "US"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-G-popular.mp3",
          "id": 37,
          "variety": "AUS"
        }
      ],
      "text": "popular"
    },
    {
      "id": 413482,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-A-train.mp3",
          "id": 55,
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-G-train.mp3",
          "id": 58,
          "variety": "US"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-B-train.mp3",
          "id": 52,
          "variety": "UK"
        }
      ],
      "text": "train"
    }
  ]
}`

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

func TestReformatJsonWithSortArrays(t *testing.T) {
	input := `{"id":"DA8GpfGvqqNf3nUQMfcCwA","words":[{"id":377534,"text":"spelling","pronunciationTracks":[{"id":44,"variety":"UK","filepath":"/pronunciations/en-GB-Standard-A-spelling.mp3"},{"id":43,"variety":"US","filepath":"/pronunciations/en-US-Standard-G-spelling.mp3"},{"id":41,"variety":"AUS","filepath":"/pronunciations/en-AUS-Standard-B-spelling.mp3"}]},{"id":304903,"text":"popular","pronunciationTracks":[{"id":32,"variety":"US","filepath":"/pronunciations/en-US-Standard-A-popular.mp3"},{"id":30,"variety":"US","filepath":"/pronunciations/en-AUS-Standard-B-popular.mp3"},{"id":37,"variety":"AUS","filepath":"/pronunciations/en-GB-Standard-G-popular.mp3"}]},{"id":413482,"text":"train","pronunciationTracks":[{"id":55,"variety":"AUS","filepath":"/pronunciations/en-AUS-Standard-A-train.mp3"},{"id":58,"variety":"US","filepath":"/pronunciations/en-US-Standard-G-train.mp3"},{"id":52,"variety":"UK","filepath":"/pronunciations/en-GB-Standard-B-train.mp3"}]}],"timestamp":"2024-06-05T22:51:09.464637Z"}`

	actual := json.ReformatJson(input, true, []string{})
	expected := `{
  "id": "DA8GpfGvqqNf3nUQMfcCwA",
  "timestamp": "2024-06-05T22:51:09.464637Z",
  "words": [
    {
      "id": 304903,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-B-popular.mp3",
          "id": 30,
          "variety": "US"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-G-popular.mp3",
          "id": 37,
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-A-popular.mp3",
          "id": 32,
          "variety": "US"
        }
      ],
      "text": "popular"
    },
    {
      "id": 377534,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-B-spelling.mp3",
          "id": 41,
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-A-spelling.mp3",
          "id": 44,
          "variety": "UK"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-G-spelling.mp3",
          "id": 43,
          "variety": "US"
        }
      ],
      "text": "spelling"
    },
    {
      "id": 413482,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-A-train.mp3",
          "id": 55,
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-B-train.mp3",
          "id": 52,
          "variety": "UK"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-G-train.mp3",
          "id": 58,
          "variety": "US"
        }
      ],
      "text": "train"
    }
  ]
}`

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

func TestReformatJsonWithExclude(t *testing.T) {
	input := `{"id":"DA8GpfGvqqNf3nUQMfcCwA","words":[{"id":377534,"text":"spelling","pronunciationTracks":[{"id":44,"variety":"UK","filepath":"/pronunciations/en-GB-Standard-A-spelling.mp3"},{"id":43,"variety":"US","filepath":"/pronunciations/en-US-Standard-G-spelling.mp3"},{"id":41,"variety":"AUS","filepath":"/pronunciations/en-AUS-Standard-B-spelling.mp3"}]},{"id":304903,"text":"popular","pronunciationTracks":[{"id":32,"variety":"US","filepath":"/pronunciations/en-US-Standard-A-popular.mp3"},{"id":30,"variety":"US","filepath":"/pronunciations/en-AUS-Standard-B-popular.mp3"},{"id":37,"variety":"AUS","filepath":"/pronunciations/en-GB-Standard-G-popular.mp3"}]},{"id":413482,"text":"train","pronunciationTracks":[{"id":55,"variety":"AUS","filepath":"/pronunciations/en-AUS-Standard-A-train.mp3"},{"id":58,"variety":"US","filepath":"/pronunciations/en-US-Standard-G-train.mp3"},{"id":52,"variety":"UK","filepath":"/pronunciations/en-GB-Standard-B-train.mp3"}]}],"timestamp":"2024-06-05T22:51:09.464637Z"}`

	actual := json.ReformatJson(input, true, []string{"words[*].pronunciationTracks[*].id"})
	expected := `{
  "id": "DA8GpfGvqqNf3nUQMfcCwA",
  "timestamp": "2024-06-05T22:51:09.464637Z",
  "words": [
    {
      "id": 304903,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-B-popular.mp3",
          "variety": "US"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-G-popular.mp3",
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-A-popular.mp3",
          "variety": "US"
        }
      ],
      "text": "popular"
    },
    {
      "id": 377534,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-B-spelling.mp3",
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-A-spelling.mp3",
          "variety": "UK"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-G-spelling.mp3",
          "variety": "US"
        }
      ],
      "text": "spelling"
    },
    {
      "id": 413482,
      "pronunciationTracks": [
        {
          "filepath": "/pronunciations/en-AUS-Standard-A-train.mp3",
          "variety": "AUS"
        },
        {
          "filepath": "/pronunciations/en-GB-Standard-B-train.mp3",
          "variety": "UK"
        },
        {
          "filepath": "/pronunciations/en-US-Standard-G-train.mp3",
          "variety": "US"
        }
      ],
      "text": "train"
    }
  ]
}`

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}
