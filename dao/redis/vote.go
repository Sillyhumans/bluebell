package redis

func GetVote(key string) (vote int8, err error) {
	res, err := client.Get(Prefix + isPostPrefix + key).Int()
	return int8(res), err
}

func SetVote(key string, vote int8) (err error) {
	_, err = client.Set(Prefix+isPostPrefix+key, vote, 0).Result()
	return err
}

func DelSetVote(key string, vote int8) (err error) {
	_, err = GetVote(key)
	if err != nil {
		err = SetVote(key, vote)
		return err
	}
	_, err = client.Do("DEL", Prefix+isPostPrefix+key).Result()
	return
}
