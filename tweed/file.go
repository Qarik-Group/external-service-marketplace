package tweed

/*func File(args []string) {
	GonnaNeedATweed()
	id, name := GonnaNeedAnInstanceAndAFile(args)

	var out []api.FileResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/instances/"+id+"/files", &out)

	if opts.JSON {
		JSON(out)
		os.Exit(0)
	}

	for _, file := range out {
		if file.Filename == name {
			fmt.Printf("%s", file.Contents)
			os.Exit(0)
		}
	}
	fmt.Fprintf(os.Stderr, "@R{%s}: file not found\n", name)
}*/
