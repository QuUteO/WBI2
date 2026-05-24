package main

func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		return nil
	}

	done := make(chan interface{})

	// Запускаем горутину для каждого канала
	for _, ch := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				// Первый закрывшийся канал закроет done
				close(done)
			case <-done:
				// Уже закрыт, выходим
			}
		}(ch)
	}

	return done
}
