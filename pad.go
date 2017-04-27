package tui

import (
  "github.com/nasciiboy/morg/katana"
)

type Pad struct {
  Buffer [][]Cell
	Curs   Gps
  Frame  Gps

  Color  uint8
  Attrs  uint8

  AutoFill bool

  Screen *Window
}

func NewPad( w *Window ) *Pad {
  return &Pad{
    Buffer: make([][]Cell, w.Height ),
    Screen: w,
  }
}

func (p *Pad) AddCh( ch uint64 ) {
  p.AddCell( extractCell( ch ) )
}

func (p *Pad) AddChs( ch []uint64 ) {
  for _, c := range ch {
    p.AddCh( c )
  }
}

func (p *Pad) AddRune( r rune ){
  c := p.GetDefaultCell()
  c.Ch = r
  p.AddCell( c )
}

func (p *Pad) AddStr( str string ){
  p.AddCells( p.StrToCells( str ) )
}

func (p *Pad) AddCell( cell Cell ) {
  cell.Touch = true

  p.Mv( p.Curs.Y, p.Curs.X )
  p.Buffer[ p.Curs.Y ][ p.Curs.X ] = cell
  p.mvCurs( cell.Ch == '\n' )
}

func (p *Pad) AddCells( cells []Cell ) {
  for _, c := range cells {
    p.AddCell( c )
  }
}

func (p *Pad) SetCell( cell Cell ) {
  cell.Touch = true

  p.Mv( p.Curs.Y, p.Curs.X )
  p.Buffer[ p.Curs.Y ][ p.Curs.X ] = cell
}


func (p *Pad) SetCells( cells []Cell ) {
  for _, c := range cells {
    p.SetCell( c )
    p.Curs.X++
  }
}

func (p *Pad) SetFace( face uint64 ){
  p.Attrs, p.Color, _, _ = extractData( face )
}

func (p *Pad) GetDefaultCell() Cell {
  return Cell{ Color: p.Color, Attrs: p.Attrs, Touch: true }
}

func (p *Pad) StrToCells( str string ) []Cell {
  c := make([]Cell, len( str ) )
  i, dCell := 0, p.GetDefaultCell()


  for _, ch := range str {
    c[i]    = dCell
    c[i].Ch = ch
    i++
  }

  return c[:i]
}

func (p *Pad) mvCurs( nl bool ) {
  if nl || (p.AutoFill && p.Curs.X + 1 >= p.Screen.Width) {
    p.Curs.X = 0
  } else {
    p.Curs.X++
    return
  }

  p.Curs.Y++
}

func (p *Pad) Mv( y, x int ){
  if y < 0 || x < 0 { return }

  for len(p.Buffer) < y + 1 {
    p.Buffer = append( p.Buffer, []Cell{} )
  }

  for len(p.Buffer[y]) < x + 1  {
    p.Buffer[y] = append(p.Buffer[y], Cell{} )
  }

  p.Curs.Y, p.Curs.X = y, x
}

func (p *Pad) Draw() {
  dCell := p.Screen.GetDefaultCell()

  for row := 0; row < p.Screen.Height; row++ {
    for col := 0; col < p.Screen.Width; col++ {
      bRow, bCol := row + p.Frame.Y, col + p.Frame.X
      if bRow < 0 || bCol < 0 || bRow >= len(p.Buffer) || bCol >= len(p.Buffer[bRow]) {
        p.Screen.Buffer[row][col] = dCell
        continue
      }

      p.Screen.Buffer[row][col] = p.Buffer[bRow][bCol]
    }
  }

  p.Screen.Draw()
}

const ( Right int = iota; Up; Left; Down; DownRight; DownLeft; UpRight;  UpLeft; PgUp; PgDown; Start; End )

func (p *Pad) Scroll( dir int ){
  switch dir {
  case Right    : p.Frame.X++
  case Up       : p.Frame.Y--
  case Left     : p.Frame.X--
  case Down     : p.Frame.Y++
  case DownRight: p.Frame.X++;  p.Frame.Y++
  case UpRight  : p.Frame.X++;  p.Frame.Y--
  case UpLeft   : p.Frame.X--;  p.Frame.Y--
  case DownLeft : p.Frame.X--;  p.Frame.Y++
  case Start    : p.Frame.X  =  0; p.Frame.Y = 0
  case End      : p.Frame.X  =  0; p.Frame.Y = len(p.Buffer) - p.Screen.Height
  case PgUp     : p.Frame.Y += -p.Screen.Height
  case PgDown   : p.Frame.Y +=  p.Screen.Height
  }

  if p.Frame.X < 0 { p.Frame.X = 0
  } else if p.Frame.X > p.Screen.Width { p.Frame.X = p.Screen.Width }

  if p.Frame.Y < 0 || len(p.Buffer) - p.Screen.Height < 0 {
    p.Frame.Y = 0
  } else if p.Frame.Y > len(p.Buffer) - p.Screen.Height {
    p.Frame.Y = len(p.Buffer) - p.Screen.Height
  }

  p.Draw()
}

func (p *Pad) AddCenterCells( cells []Cell, leftMargin uint ) {
  width  := p.Screen.Width - int(leftMargin)
  margin := make( []Cell, p.Screen.Width )

  for i, w := 0, 0; i < len(cells); i += w {
    if cells[i].Ch == ' ' {
      w = 1
      continue
    }

    w = len( cells[i:] )

    if w > width {
      if cells[width].Ch == ' ' {
        w = width
      } else {
        for s := width; s > 0; s-- {
          if cells[i + s].Ch == ' ' {
            w = s
            break
          }
        }
      }
    }

    ow := uint((width - w) / 2 ) + leftMargin
    p.SetCells( margin[:ow] )
    p.SetCells( cells[i:i + w] )
    p.mvCurs( true )
  }

  p.mvCurs( true )
}

func (p *Pad) AddLeftCells( cells []Cell, width, leftMargin uint ) {
  margin := make( []Cell, leftMargin )

  for i, w := 0, 0; i < len(cells); i += w {
    if cells[i].Ch == ' ' {
      w = 1
      continue
    }

    p.SetCells( margin )

    w = len( cells[i:] )

    if w > int(width) {
      if cells[width].Ch == ' ' {
        w = int(width)
      } else {
        for s := int(width); s > 0; s-- {
          if cells[i + s].Ch == ' ' {
            w = s
            break
          }
        }
      }
    }

    p.SetCells( cells[i:i + w] )
    p.mvCurs( true )
  }
}

func (p *Pad) AddRightCells( cells []Cell, leftMargin, rightMargin uint ) {
  lMargin := make( []Cell, leftMargin )
  rMargin := make( []Cell, p.Screen.Width - int(leftMargin) )
  width   := p.Screen.Width - int(leftMargin) - int(rightMargin)

  for i, w, fill := 0, 0, 0; i < len(cells); i += w {
    if cells[i].Ch == ' ' {
      w = 1
      continue
    }

    p.SetCells( lMargin )
    w = len( cells[i:] )

    if w > int(width) {
      if cells[width].Ch == ' ' {
        w = int(width)
      } else {
        for s := int(width); s > 0; s-- {
          if cells[i + s].Ch == ' ' {
            w = s
            break
          }
        }
      }
    }
    fill = width - w
    p.SetCells( rMargin[:fill] )
    p.SetCells( cells[i:i + w] )
    p.mvCurs( true )
  }
}

func (p *Pad) AddPreCells( cells []Cell, leftMargin uint ) {
  margin := make( []Cell, leftMargin )
  p.AddCells( margin )
  for i := 0; i < len( cells ); i++ {
    if cells[i].Ch == '\n' {
      p.mvCurs( true )
      p.SetCells( margin )
      continue
    }

    p.SetCell( cells[i] )
    p.Curs.X++
  }

  p.mvCurs( true )
}

func (p *Pad) ParseMorg( str string ) {
  var doc katana.Doc
  doc.Parse( str )

  p.Mv( 2, 0 )
  p.AddCenterCells( fontify( doc.Title, 0 ), 0 )

  if doc.Subtitle.HasSomething() {
    p.AddCenterCells( fontify( doc.Subtitle, 0 ), 0 )
    p.mvCurs( true )
  }

  p.mvCurs( true )

  p.makeBody( doc.Toc, doc.OptionsData )
}

func (p *Pad) makeBody( toc []katana.DocNode, options katana.Options ) {
  for _, h := range toc {
    full := h.Get()
    if full.N == 0 {
    } else {
      for i := 0; i < full.N; i++ {
        p.AddCh( '#' | ColorCyan | Bold )
      }

      p.AddCh( ' ' | ColorCyan | Bold )

      p.CustomFontify( full.Mark, ColorCyan | Bold )
      p.mvCurs( true )
    }

    if len( h.Cont ) > 0 {
      p.mvCurs( true )
      p.walkContent( h.Cont, options, 0 )
    }
  }
}

func (p *Pad) walkContent( doc []katana.DocNode, options katana.Options, deep uint ){
  for _, node := range doc {
    switch node.Type() {
    case katana.EmptyNode     :
    case katana.CommentNode   :
    case katana.CommandNode   : p.makeCommand( node, options, deep + 2 )
    case katana.HeadlineNode  :
  //   case katana.TableNode     : str += makeTable  ( full, c.Cont, options )
    case katana.ListNode      : p.makeList   ( node, options, deep + 2 )
    case katana.AboutNode     : p.makeAbout  ( node, options, deep + 2 )
    case katana.TextNode      :
      marginLeft := deep + 2
      width      := uint(p.Screen.Width) - 2 -10 - deep
      data := node.Get()
      p.AddLeftCells( CustomFontify( data.Mark, ColorBW ), width, marginLeft )
      p.mvCurs( true )
    }
  }
}


func (p *Pad) Shoot( y, x int, cell Cell ){
  for len(p.Buffer) < y + 1 {
    p.Buffer = append( p.Buffer, []Cell{} )
  }

  for len(p.Buffer[y]) < x + 1  {
    p.Buffer[y] = append(p.Buffer[y], Cell{} )
  }

  p.Buffer[y][x] = cell
}

func (p *Pad) Shooter( y, x int, cells []Cell ){
  for _, cell := range cells {
    p.Shoot( y, x, cell )
    x++
  }
}

func (p *Pad) makeList( node katana.DocNode, options katana.Options, deep uint ){
  dNode := node.Get()

  for _, element := range node.Cont {
    full := element.Get()
    row := p.Curs.Y
    switch dNode.N {
    case katana.ListMinusNode, katana.ListPlusNode :
      p.walkContent( element.Cont, options, deep )
      p.Shoot( row, int(deep), Cell{ Ch:  '-', Color: cBG, Attrs: aBold } )
    case katana.ListNumNode,   katana.ListAlphaNode:
      str := full.Data + "."
      p.walkContent( element.Cont, options, deep + uint(len(str)) - 1 )
      p.Shooter( row, int(deep), CustomFontify( katana.StrToMark( str ), ColorBG | Bold ) )
    case katana.ListMdefNode,  katana.ListPdefNode :
      p.makeDlListNode( element, options, deep )
    case katana.ListDialogNode:
      p.walkContent( element.Cont, options, deep )
      p.Shoot( row, int(deep), Cell{ Ch:  '>', Color: cBG, Attrs: aBold } )
    }
  }


  return
}

func (p *Pad) makeDlListNode( node katana.DocNode, options katana.Options, deep uint ){
  data  := node.Get()
  width := uint(p.Screen.Width) - 1 - 10 - deep
  p.AddLeftCells( CustomFontify( data.Mark, ColorRed | Bold ), width, deep + 1 )

  p.walkContent( node.Cont, options, deep + 3 )
}

func (p *Pad) makeAbout( node katana.DocNode, options katana.Options, deep uint ){
  data := node.Get()

  width := uint(p.Screen.Width) - 1 - 10 - deep
  p.AddLeftCells( CustomFontify( data.Mark, ColorBW | Bold ), width, deep + 2 )
  p.mvCurs( true )

  p.walkContent( node.Cont, options, deep + 2 )
}

func (p *Pad) makeCommand( node katana.DocNode, options katana.Options, deep uint ){
  data := node.Get()

  switch data.Comm {
  case "src"    : p.makeCommandSrc( data, options, deep )
  case "figure" : p.makeCommandFigure( data, node.Cont, options, deep )
  // case "cols"   : return makeCommandCols  ( data, node.Cont, options )
  // case "img"    : return makeCommandImg   ( data, node.Cont, options )
  // case "video"  : return makeCommandVideo ( data, node.Cont, options )
  case "quote"  : p.makeCommandQuote ( data, node.Cont, options, deep )
  case "example", "pre", "diagram", "art":
    p.makeCommandPre( data, options, deep )
  case "center", "bold", "emph", "italic":
    p.makeCommandFont( data, node.Cont, options, deep )
  }
}

func (p *Pad) makeCommandSrc( comm katana.FullData, options katana.Options, deep uint ){
  p.AddPreCells( StrToCells( comm.Data, cBW, aBold ), deep + 4 )
  p.mvCurs( true )
}

func (p *Pad) makeCommandPre( comm katana.FullData, options katana.Options, deep uint ){
  p.AddPreCells( StrToCells( comm.Data, cBW, aBold ), deep + 4 )
  p.mvCurs( true )
}

func (p *Pad) makeCommandFont( comm katana.FullData, body []katana.DocNode, options katana.Options, deep uint ){
  for _, node := range body {
    switch node.Type() {
    default:
      bNode := []katana.DocNode{ node }
      p.walkContent( bNode, options, deep + 2 )
    case katana.TextNode      :
      switch comm.Comm {
      case "center":
        data := node.Get()
        p.AddCenterCells( CustomFontify( data.Mark, ColorBW ), 0 )
        p.mvCurs( true )
      case "bold":
        marginLeft := deep
        width      := uint(p.Screen.Width) - 2 -10 - deep
        data := node.Get()
        p.AddLeftCells( CustomFontify( data.Mark, ColorBW | Bold ), width, marginLeft )
        p.mvCurs( true )
      case "emph":
        marginLeft := deep
        width      := uint(p.Screen.Width) - 2 -10 - deep
        data := node.Get()
        p.AddLeftCells( CustomFontify( data.Mark, ColorBY | Bold ), width, marginLeft )
        p.mvCurs( true )
      case "italic":
      }
    }
  }
}

func (p *Pad) makeCommandQuote( comm katana.FullData, body []katana.DocNode, options katana.Options, deep uint ){
  for _, c := range body {
    nodeData := c.Get()
    switch nodeData.N {
    case katana.TextQuoteAuthor:
      p.AddRightCells( CustomFontify( nodeData.Mark, ColorBM | Bold ), deep, 10 )
    case katana.TextSimple:
      marginLeft := deep + 2
      width      := uint(p.Screen.Width) - 2 -10 - deep
      p.AddLeftCells( CustomFontify( nodeData.Mark, ColorBY | Bold ), width, marginLeft )
      p.mvCurs( true )
    }
  }

  p.mvCurs( true )
}

func (p *Pad) makeCommandFigure( comm katana.FullData, body []katana.DocNode, options katana.Options, deep uint ){
  marginLeft := deep + 2
  width      := uint(p.Screen.Width) - 2 -10 - deep
  p.AddLeftCells( CustomFontify( comm.Mark, ColorBC | Bold ), width, marginLeft )
  p.mvCurs( true )
  p.walkContent( body, options, deep + 3 )
}

func fontify( m katana.Markup, attrs uint8 ) []Cell {
  att, color, _, _ := extractData( getColor( m.Type ) )

  if len( m.Custom ) == 0 && len( m.Body ) == 0 {
    return StrToCells( m.Data, color, att | attrs )
  }

  var body CellBuffer
  for _, c := range m.Body {
    if c.Type == katana.MarkupNil {
      c.Type = m.Type
    } else {
      att = extractAttrs( getColor( m.Type ) )
    }

    body.Write( fontify( c, attrs | att ) )
  }

  return atCommand( body.Data(), m.Type, attrs )
}

func (p *Pad) CustomFontify( m katana.Markup, color uint64 ){
  if len( m.Custom ) == 0 && len( m.Body ) == 0 {
    for _, ch := range m.Data {
      p.AddCh( uint64(ch) | color )
    }
    return
  }

  for _, c := range m.Body {
    if c.Type == katana.MarkupNil {
      p.CustomFontify( c, color )
    } else {
      p.CustomFontify( c, getColor( c.Type ) | getAttrs( color ))
    }
  }
}

// func tontify( m katana.Markup ) (str string) {
//   if len( m.Custom ) == 0 && len( m.Body ) == 0 {
//     return m.Data
//   }

//   var custom, body string
//   for _, c := range m.Custom { custom += tontify( c ) }
//   for _, c := range m.Body   { body   += tontify( c ) }

//   if custom == "" {
//     switch m.Type {
//     case 'l', 'N', 'n', 't' :  custom  = m.MakeCustom()
//     }
//   }

//   return atCommand( body, custom, m.Type )
// }

func CustomFontify( m katana.Markup, color uint64 ) []Cell {
  var body CellBuffer

  if len( m.Custom ) == 0 && len( m.Body ) == 0 {
    for _, ch := range m.Data {
      body.WriteU64( uint64(ch) | color )
    }

    return body.Data()
  }

  for _, c := range m.Body {
    if c.Type == katana.MarkupNil {
      body.Write( CustomFontify( c, color ) )
    } else {
      // body.Write( openDecorator( c.Type ) )
      body.Write( CustomFontify( c, getColor( c.Type ) | getAttrs( color ) ) )
      // body.Write( closeDecorator( c.Type ) )
    }
  }

  return body.Data()
}

func Fontify( str string ) []Cell {
  var markup katana.Markup
  markup.Parse( str )

  return fontify( markup, 0 )
}

func ToText( str string ) []uint64 {
  var markup katana.Markup
  markup.Parse( str )

  return ToCustomU64( markup.String(), 0 )
}

func ToCustomU64( str string, color uint64 ) (result []uint64) {
  // result = make( []uint64, 0, 32 )

  // for _, c := range str {
  //   result = append( result, uint64( c ) | color )
  // }

  return
}

func getColor( t byte ) uint64 {
  switch t {
  case katana.MarkupNil, katana.MarkupEsc, katana.MarkupErr: return ColorBW
  case katana.MarkupHeadline: return ColorCyan | Bold
  case katana.MarkupText    : return ColorWhite
  case katana.MarkupTitle   : return ColorBW | Bold
  case katana.MarkupSubTitle: return ColorBR | Bold
  case '!' : return ColorWhite
  case '"' : return ColorWhite
  case '#' : return ColorWhite
  case '$' : return ColorRed | Bold
  case '%' : return ColorWhite
  case '&' : return ColorWhite
  case '\'': return ColorWhite
  case '*' : return ColorWhite
  case '+' : return ColorWhite
  case ',' : return ColorWhite
  case '-' : return ColorWhite
  case '.' : return ColorWhite
  case '/' : return ColorWhite
  case '0' : return ColorWhite
  case '1' : return ColorWhite
  case '2' : return ColorWhite
  case '3' : return ColorWhite
  case '4' : return ColorWhite
  case '5' : return ColorWhite
  case '6' : return ColorWhite
  case '7' : return ColorWhite
  case '8' : return ColorWhite
  case '9' : return ColorWhite
  case ':' : return ColorWhite
  case ';' : return ColorWhite
  case '=' : return ColorWhite
  case '?' : return ColorWhite
  case 'A' : return ColorWhite
  case 'B' : return ColorWhite
  case 'C' : return ColorWhite
  case 'D' : return ColorWhite
  case 'E' : return ColorWhite
  case 'F' : return ColorWhite
  case 'G' : return ColorWhite
  case 'H' : return ColorWhite
  case 'I' : return ColorWhite
  case 'J' : return ColorWhite
  case 'K' : return ColorWhite
  case 'L' : return ColorWhite
  case 'M' : return ColorWhite
  case 'N' : return ColorWhite
  case 'O' : return ColorWhite
  case 'P' : return ColorWhite
  case 'Q' : return ColorWhite
  case 'R' : return ColorWhite
  case 'S' : return ColorWhite
  case 'T' : return ColorWhite
  case 'U' : return ColorWhite
  case 'V' : return ColorWhite
  case 'W' : return ColorWhite
  case 'X' : return ColorWhite
  case 'Y' : return ColorWhite
  case 'Z' : return ColorWhite
  case '\\': return ColorWhite
  case '^' : return ColorWhite
  case '_' : return ColorWhite
  case '`' : return ColorWhite
  case 'a' : return ColorWhite
  case 'b' : return ColorWhite | Bold
  case 'c' : return ColorBlue  | Bold
  case 'd' : return ColorWhite
  case 'e' : return ColorMagenta | Bold
  case 'f' : return ColorWhite
  case 'g' : return ColorWhite
  case 'h' : return ColorWhite
  case 'i' : return ColorWhite
  case 'j' : return ColorWhite
  case 'k' : return ColorWhite
  case 'l' : return ColorYellow
  case 'm' : return ColorWhite
  case 'n' : return ColorWhite
  case 'o' : return ColorWhite
  case 'p' : return ColorWhite
  case 'q' : return ColorWhite
  case 'r' : return ColorWhite
  case 's' : return ColorWhite
  case 't' : return ColorWhite
  case 'u' : return ColorWhite
  case 'v' : return ColorWhite
  case 'w' : return ColorWhite
  case 'x' : return ColorWhite
  case 'y' : return ColorWhite
  case 'z' : return ColorWhite
  case '|' : return ColorWhite
  case '~' : return ColorWhite
  }

  return ColorRed
}

func atCommand( body []Cell, t, attrs uint8 ) []Cell {
  var buff CellBuffer
  buff.SetFace( getColor( t ) | uint64(attrs) << 56 )

  switch t {
  case katana.MarkupNil, katana.MarkupEsc, katana.MarkupErr: return body
  case katana.MarkupHeadline: return body
  case katana.MarkupText    : return body
  case katana.MarkupTitle   : return body
  case katana.MarkupSubTitle: return body
  case '!' : return body
  case '"' :
    buff.WriteRune( '“' )
    buff.Write( body )
    buff.WriteRune( '”' )
  case '#' : return body
  case '$' : return body
  case '%' : return body
  case '&' : return body
  case '\'':
    buff.WriteRune( '‘' )
    buff.Write( body )
    buff.WriteRune( '’' )
  case '*' : return body
  case '+' : return body
  case ',' : return body
  case '-' :
    buff.WriteString( "––" )
    buff.Write( body )
    buff.WriteString( "––" )
  case '.' : return body
  case '/' : return body
  case '0' : return body
  case '1' : return body
  case '2' : return body
  case '3' : return body
  case '4' : return body
  case '5' : return body
  case '6' : return body
  case '7' : return body
  case '8' : return body
  case '9' : return body
  case ':' : return body
  case ';' : return body
  case '=' : return body
  case '?' : return body
  case 'A' : return body
  case 'B' : return body
  case 'C' : return body
  case 'D' : return body
  case 'E' : return body
  case 'F' : return body
  case 'G' : return body
  case 'H' : return body
  case 'I' : return body
  case 'J' : return body
  case 'K' : return body
  case 'L' : return body
  case 'M' : return body
  case 'N' : return body
  case 'O' : return body
  case 'P' : return body
  case 'Q' : return body
  case 'R' : return body
  case 'S' : return body
  case 'T' : return body
  case 'U' : return body
  case 'V' : return body
  case 'W' : return body
  case 'X' : return body
  case 'Y' : return body
  case 'Z' : return body
  case '\\': return body
  case '^' : return body
  case '_' : return body
  case '`' : return body
  case 'a' : return body
  case 'b' : return body
  case 'c' : return body
  case 'd' : return body
  case 'e' : return body
  case 'f' : return body
  case 'g' : return body
  case 'h' : return body
  case 'i' : return body
  case 'j' : return body
  case 'k' : return body
  case 'l' : return body
  case 'm' : return body
  case 'n' : return body
  case 'o' : return body
  case 'p' : return body
  case 'q' :
    buff.WriteRune( '“' )
    buff.Write( body )
    buff.WriteRune( '”' )
  case 'r' : return body
  case 's' : return body
  case 't' : return body
  case 'u' : return body
  case 'v' : return body
  case 'w' : return body
  case 'x' : return body
  case 'y' : return body
  case 'z' : return body
  case '|' : return body
  case '~' : return body
  }

  return buff.Data()
}

func openDecorator( t uint8 ) string {
  switch t {
  case '"', 'q' : return "“"
  case '\'': return "‘"
  case '-' : return "––"
  }

  return ""
}

func closeDecorator( t uint8 ) string {
  switch t {
  case '"', 'q' : return "”"
  case '\'': return "’"
  case '-' : return "––"
  }

  return ""
}
