package token

func readSplitter(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var start int
	for start < len(data) && isSpace(data[start]) {
		start++
	}

	for idx := start; idx < len(data); idx++ {
		switch data[idx] {
		case '(', ')', '{', '}':
			if idx == start {
				return idx + 1, data[start : idx+1], nil
			}
			return idx, data[start:idx], nil

		case ' ', '\r', '\n', '\t':
			return idx + 1, data[start:idx], nil
		}
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}

func isSpace(space byte) bool {
	switch space {
	case ' ', '\r', '\n', '\t':
		return true
	default:
		return false
	}
}
