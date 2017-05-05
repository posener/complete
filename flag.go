package complete

type FlagOptions struct {
	HasFollow      bool
	FollowsOptions []string
}

var (
	FlagNoFollow      = FlagOptions{}
	FlagUnknownFollow = FlagOptions{HasFollow: true}
)

