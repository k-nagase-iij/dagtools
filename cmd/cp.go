package cmd

import (
	"flag"
	"fmt"
	"github.com/iij/dagtools/client"
	"github.com/iij/dagtools/env"
	"os"
	"strings"
)

type cpCommand struct {
	env   *env.Environment
	cli   client.StorageClient
	opts  *flag.FlagSet
	force bool
}

func (c *cpCommand) Description() string {
	return "copy object or directory"
}

func (c *cpCommand) Usage() string {
	return fmt.Sprintf(`Command Usage:
  cp [-f] [<bucket>]:[<key>] <bucket>
  cp [-f] [<bucket>]:[<key>] <bucket>:<prefix>

Options:
%s`, OptionUsage(c.opts))
}

func (c *cpCommand) Init(env *env.Environment) (err error) {
	c.env = env
	c.cli, _ = client.NewStorageClient(env)
	opts := flag.NewFlagSet("cp", flag.ExitOnError)
	opts.BoolVar(&c.force, "f", false, "Ignore modifies")
	opts.Usage = func() {
		fmt.Fprintln(os.Stdout, c.Usage())
	}
	c.opts = opts
	return
}

func (c *cpCommand) Run(args []string) (err error) {
	var (
		sourceMeta   *client.ObjectSummary
		regions      *client.Regions
		sourceBucket = ""
		sourceKey    = ""
		destBucket   = ""
		destKey      = ""
	)
	if len(args) == 0 {
		return ErrArgument
	}
	c.opts.Parse(args)
	argv := c.opts.Args()
	if len(argv) < 2 {
		return ErrArgument
	}
	source := strings.Split(argv[0], ":")
	sourceBucket = source[0]
	sourceKey = source[1]
	if !strings.HasSuffix(argv[1], ":") {
		destBucket = argv[1]
	} else {
		dest := strings.Split(argv[1], ":")
		destBucket = dest[0]
		destKey = dest[1]
	}
	if sourceBucket == "" || sourceKey == "" || destBucket == "" {
		return ErrArgument
	}
	if destKey == "" {
		destKey = sourceKey
	}
	if c.env.Verbose {
		fmt.Fprintf(os.Stdout, "copy: %s:%s -> %s:%s\n", sourceBucket, sourceKey, destBucket, destKey)
	}
	if sourceMeta, err = c.cli.GetObjectSummary(sourceBucket, sourceKey); err != nil {
		regions, err = c.cli.GetRegions()
		sourceRegion, err := c.cli.GetBucketLocation(sourceBucket)
		if err != nil {
			return err
		}
		defaultEp := c.cli.GetEndpoint()
		for _, r := range regions.Regions {
			if r.Name == sourceRegion {
				c.cli.SetEndpoint(r.Endpoint)
				break
			}
		}
		sourceMeta, err = c.cli.GetObjectSummary(sourceBucket, sourceKey)
		if err != nil {
			return err
		}
		c.cli.SetEndpoint(defaultEp)
	}
	if err = c.cli.PutObjectCopy(sourceBucket, sourceKey, destBucket, destKey, sourceMeta); err != nil {
		regions, err = c.cli.GetRegions()
		destRegion, err := c.cli.GetBucketLocation(destBucket)
		if err != nil {
			return err
		}
		for _, r := range regions.Regions {
			if r.Name == destRegion {
				c.cli.SetEndpoint(r.Endpoint)
				break
			}
		}
		err = c.cli.PutObjectCopy(sourceBucket, sourceKey, destBucket, destKey, sourceMeta)
		if err != nil {
			return err
		}
	}
	return
}

func init() {
	Commands.Register(new(cpCommand), "cp")
}
