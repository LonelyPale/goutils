// Created by LonelyPale at 2019-12-06

package mongodb

/*
func Client(opts ...interface{}) (cli *mongo.Client, err error) {
	//异常处理必须放在抛出异常的语句之前，否则是无法处理异常的，要处理全局异常，必须放在函数的开始位置
	//panic异常处理机制不会自动将错误信息传递给error，所以要在func函数中进行显式的传递
	//必须要先声明defer，否则不能捕获到panic异常
	defer func() {
		if p := recover(); p != nil {
			log.Println("panic recover! p:", p)
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
			} else {
				err = errors.New("panic")
			}
			debug.PrintStack()
		}
	}()

	if client == nil {
		var configPath string
		if len(opts) > 0 {
			configPath = opts[0].(string)
		}

		conf, err := ReadConfig(configPath)
		if err != nil {
			log.Println(err)
			return client, err
		}
		if conf.URI == "" {
			err = errors.New("conf.URI cannot be empty")
			log.Println(err)
			return client, err
		}

		client, err = mongo.NewClient(options.Client().ApplyURI(conf.URI))
		if err != nil {
			log.Println(err)
			return client, err
		}

		timeout := time.Duration(conf.Timeout)
		ctx, _ = context.WithTimeout(context.Background(), timeout*time.Second)

		err = client.Connect(ctx)
		if err != nil {
			log.Println(err)
			return client, err
		}
	}

	return client, err
}

*/
