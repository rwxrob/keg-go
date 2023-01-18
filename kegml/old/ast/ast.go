package ast

const (
	Unknown = iota
	EndBlock
	NodeBlocks
	Title
	IncludesB
	Separator
	BulletedB
	NumberedB
	FigureB
	ShortQuoteB
	LongQuoteB
	MathB
	FencedB
	DivisionB
	ParagraphB
	FootnotesB
	Node
	Includes
	Bulleted
	Numbered
	Figure
	ShortQuote
	LongQuote
	Math
	Fenced
	Division
	Paragraph
	Footnotes
)

var Rules = []string{
	`Unknown`,
	`endBlock`,
	`NodeBlocks`,
	`Title`,
	`IncludesB`,
	`Separator`,
	`BulletedB`,
	`NumberedB`,
	`FigureB`,
	`ShortQuoteB`,
	`LongQuoteB`,
	`MathB`,
	`FencedB`,
	`DivisionB`,
	`ParagraphB`,
	`FootnotesB`,
	`Node`,
	`Includes`,
	`Bulleted`,
	`Numbered`,
	`Figure`,
	`ShortQuote`,
	`LongQuote`,
	`Math`,
	`Fenced`,
	`Division`,
	`Paragraph`,
	`Footnotes`,
}
