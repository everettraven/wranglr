package tabs

type Options func(*Model)

func WithStyles(styles Styles) Options {
	return func(t *Model) {
		t.styles = styles
	}
}

func WithDisplayFormat(display DisplayFormat) Options {
	return func(t *Model) {
		t.display = display

		if display == DisplayFormatHorizontal {
			t.keys = DefaultHorizontalKeyMap
			t.styles = DefaultHorizontalStyles
		}

		if t.display == DisplayFormatVertical {
			t.keys = DefaultVerticalKeyMap
			t.styles = DefaultVerticalStyles
		}
	}
}

func WithKeyMap(keys KeyMap) Options {
	return func(t *Model) {
		t.keys = keys
	}
}

func WithWrapping(wrap bool) Options {
	return func(t *Model) {
		t.wrap = wrap
	}
}

func WithVerticalWidth(width int) Options {
	return func(t *Model) {
		t.width = width
	}
}
