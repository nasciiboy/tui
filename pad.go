package tui

import (
  "github.com/nasciiboy/morg/katana"
)

type Pad struct {
  Buffer [][]Cell
	Curs   Gps
  Frame  Gps

  Screen *Window
}

func NewPad( w *Window ) *Pad {
  return &Pad{
    Buffer: make([][]Cell, w.Height ),
    Screen: w,
  }
}

func (p *Pad) AddCh( ch uint64 ) {
  chAttrs, chColor, _, r := extractData( ch )
  cell := Cell{ Attrs: chAttrs, Color: chColor, Ch: r, Touch: true }

  p.AddCell( cell )
}

func (p *Pad) AddChs( ch []uint64 ) {
  for _, c := range ch {
    p.AddCh( c )
  }
}

func (p *Pad) AddStr( str string ){
  for _, c := range( str ) {
    p.AddCh( uint64(c) )
  }
}

func (p *Pad) AddCell( cell Cell ) {
  p.Mv( p.Curs.Y, p.Curs.X )
  p.Buffer[ p.Curs.Y ][ p.Curs.X ] = cell
  p.mvCurs( cell.Ch == '\n' )
}

func (p *Pad) SetCell( cell Cell ) {
  p.Mv( p.Curs.Y, p.Curs.X )
  p.Buffer[ p.Curs.Y ][ p.Curs.X ] = cell
}

func (p *Pad) AddCells( cells []Cell ) {
  for _, c := range cells {
    p.AddCell( c )
  }
}

func (p *Pad) SetCells( cells []Cell ) {
  for _, c := range cells {
    p.SetCell( c )
    p.Curs.X++
  }
}

func (p *Pad) mvCurs( nl bool ) {
  if nl || p.Curs.X + 1 >= p.Screen.Width {
    p.Curs.X = 0
  } else {
    p.Curs.X++
    return
  }

  if p.Curs.Y + 1 < len(p.Buffer) {
    p.Curs.Y++
  } else {
    p.Buffer = append(p.Buffer, make( []Cell, 0, 4 ) )
    p.Curs.Y++
  }
}

func (p *Pad) Draw() {
  for row := 0; row < p.Screen.Height; row++ {
    for col := 0; col < p.Screen.Width; col++ {
      bRow, bCol := row + p.Frame.Y, col + p.Frame.X
      if bRow < len(p.Buffer) && bCol < len(p.Buffer[bRow]) {
        p.Screen.Buffer[row][col] = p.Buffer[bRow][bCol]
        p.Screen.Buffer[row][col].Touch = true
        continue
      }

      p.Screen.Buffer[row][col] = Cell{ Touch: true }
    }
  }

  p.Screen.Touch = true
  p.Screen.Refresh()
}

const ( Right int = iota; Up; Left; Down; DownRight; DownLeft; UpRight;  UpLeft; PgUp; PgDown; Start; End )

func (p *Pad) Scroll( dir int ){
  var dx, dy int

  switch dir {
  case Right    : dx =  1; dy =  0
  case Up       : dx =  0; dy = -1
  case Left     : dx = -1; dy =  0
  case Down     : dx =  0; dy =  1
  case DownRight: dx =  1; dy =  1
  case UpRight  : dx =  1; dy = -1
  case UpLeft   : dx = -1; dy = -1
  case DownLeft : dx = -1; dy =  1
  case PgUp   :
    dy = -p.Screen.Height
    if p.Frame.Y + dy < p.Screen.Height {
      p.Frame.Y = 0
    }
  case PgDown :
    dy = p.Screen.Height
    if p.Frame.Y + dy > len(p.Buffer) - p.Screen.Height {
      p.Frame.Y = len(p.Buffer) - p.Screen.Height
      dy = 0
    }
  case Start: p.Frame.X = 0; p.Frame.Y = 0
  case End  : p.Frame.X = 0;
    p.Frame.Y = len(p.Buffer) - p.Screen.Height
  }

  if p.Frame.X + dx >= 0 && p.Frame.Y + dy >= 0 &&
     p.Frame.X + dx <= p.Screen.Width && p.Frame.Y + dy <= len(p.Buffer) {
    p.Frame.X += dx
    p.Frame.Y += dy
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

  p.mvCurs( true )
  p.mvCurs( true )


  p.AddCenterCells( CustomFontify( katana.StrToMark( doc.Title ), ColorCyan | Bold ), 0 )

  if len(doc.Subtitle) != 0 {
    p.AddCenterCells( CustomFontify( katana.StrToMark( doc.Subtitle ), ColorRed | Bold ), 0 )
    p.mvCurs( true )
  }


// {{ if .OptionsData.Toc }}
// <div id="toc">
//   <p>index</p>
//   <div id="toc-contents">
//   {{ ToToc .Toc .OptionsData }}
//   </div>
// </div>
// {{ end }}

  p.mvCurs( true )

  p.makeBody( doc.Toc, doc.OptionsData )
}

func (p *Pad) makeBody( toc []katana.DocNode, options katana.Options ) {
  for _, h := range( toc ) {
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
  for _, node := range( doc ) {
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

func (p *Pad) Mv( y, x int ){
  for len(p.Buffer) < y + 1 {
    p.Buffer = append( p.Buffer, []Cell{} )
  }

  for len(p.Buffer[y]) < x + 1  {
    p.Buffer[y] = append(p.Buffer[y], Cell{} )
  }

  p.Curs.Y, p.Curs.X = y, x
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

  for _, element := range( node.Cont ) {
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
  for _, node := range( body ) {
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
  for _, c := range( body ) {
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

func fontify( m katana.Markup ) (result []Cell) {
  if len( m.Custom ) == 0 && len( m.Body ) == 0 {
    result = make( []Cell, len(m.Data) )

    i := 0
    for _, ch := range m.Data {
      result[ i ].Ch = ch
      i++
    }

    return result[:i]
  }

  custom, body := make( []Cell, 0, 32 ), make( []Cell, 0, 32 )

  for _, c := range m.Custom {
    p := fontify( c )
    custom = append( custom, p... )
  }

  for _, c := range m.Body {
    p := fontify( c )
    body = append( body, p... )
  }

  if len(custom) == 0 {
    switch m.Type {
    case 'l', 'N', 'n', 't' :
      // custom = ToU64( m.MakeCustom(), m.Type )
    }
  }

//  return ToLabel( body, custom, m.Type )
  return
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
      p.CustomFontify( c, getColor( c.Type ) | extractAttrs( color ))//      body.Write( customFontify( c,  ) )
    }
  }
}

func CustomFontify( m katana.Markup, color uint64 ) []Cell {
  var body Buffer

  if len( m.Custom ) == 0 && len( m.Body ) == 0 {
    for _, ch := range m.Data {
      body.WriteU64( uint64(ch) | color )
    }

    return body.CellData()
  }

  for _, c := range m.Body {
    if c.Type == katana.MarkupNil {
      body.WriteCells( CustomFontify( c, color ) )
    } else {
      body.WriteCells( openDecorator( c.Type ) )
      body.WriteCells( CustomFontify( c, getColor( c.Type ) | extractAttrs( color ) ) )
      body.WriteCells( closeDecorator( c.Type ) )
    }
  }

  return body.CellData()
}

func Fontify( str string ) []Cell {
  var markup katana.Markup
  markup.Parse( str )

  return fontify( markup )
}

func ToText( str string ) []uint64 {
  var markup katana.Markup
  markup.Parse( str )

  return ToCustomU64( markup.String(), 0 )
}

func ToCustomU64( str string, color uint64 ) (result []uint64) {
  // result = make( []uint64, 0, 32 )

  // for _, c := range( str ) {
  //   result = append( result, uint64( c ) | color )
  // }

  return
}

func getColor( t byte ) uint64 {
  switch t {
  case katana.MarkupNil, katana.MarkupEsc, katana.MarkupErr: return ColorBW
  case katana.MarkupHeadline: return ColorCyan | Bold
  case katana.MarkupText    : return ColorWhite
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

func (p *Pad) Fontify( m katana.Markup, t uint8 ){
  switch t {
  case katana.MarkupNil, katana.MarkupEsc, katana.MarkupErr: p.CustomFontify( m, ColorBW )
  case '!' : p.CustomFontify( m, ColorBW )
  case '"' :
    p.AddCh( '“' | getColor( t ) )
    p.CustomFontify( m, getColor( '"') )
    p.AddCh( '”' | getColor( t ) )
  case '#' : p.CustomFontify( m, ColorBW ) //`<span class="path" >` + body + "</span>"
  case '$' : p.CustomFontify( m, ColorBW ) //`<code class="command" >` + body + "</code>"
  case '%' : p.CustomFontify( m, ColorBW ) //body // "parentesis"
  case '&' : p.CustomFontify( m, ColorBW ) //body
  case '\'': p.CustomFontify( m, ColorBW )
    // buf.Write( ToU64( "‘", '\'' ) )
    // buf.Write( body )
    // buf.Write( ToU64( "’", '\'' ) )
    // return buf.Data()
  case '*' : p.CustomFontify( m, ColorBW ) //body
  case '+' : p.CustomFontify( m, ColorBW ) //body
  case ',' : p.CustomFontify( m, ColorBW ) //body
  case '-' : p.CustomFontify( m, ColorBW )
    // buf.Write( ToU64( "––", '-' ) )
    // buf.Write( body )
    // buf.Write( ToU64( "––", '-' ) )
    // return buf.Data()
  case '.' : p.CustomFontify( m, ColorBW ) //body
  case '/' : p.CustomFontify( m, ColorBW ) //body
  case '0' : p.CustomFontify( m, ColorBW ) //body
  case '1' : p.CustomFontify( m, ColorBW ) //body
  case '2' : p.CustomFontify( m, ColorBW ) //body
  case '3' : p.CustomFontify( m, ColorBW ) //body
  case '4' : p.CustomFontify( m, ColorBW ) //body
  case '5' : p.CustomFontify( m, ColorBW ) //body
  case '6' : p.CustomFontify( m, ColorBW ) //body
  case '7' : p.CustomFontify( m, ColorBW ) //body
  case '8' : p.CustomFontify( m, ColorBW ) //body
  case '9' : p.CustomFontify( m, ColorBW ) //body
  case ':' : p.CustomFontify( m, ColorBW ) //"<dfn>" + body + "</dfn>"
  case ';' : p.CustomFontify( m, ColorBW ) //body
  case '=' : p.CustomFontify( m, ColorBW ) //body
  case '?' : p.CustomFontify( m, ColorBW ) //body
  case 'A' : p.CustomFontify( m, ColorBW ) //`<span class="acronym" >` + body + "</span>"
  case 'B' : p.CustomFontify( m, ColorBW ) //body
  case 'C' : p.CustomFontify( m, ColorBW ) //body // "smallCaps"
  case 'D' : p.CustomFontify( m, ColorBW ) //body
  case 'E' : p.CustomFontify( m, ColorBW ) //body // "error"
  case 'F' : p.CustomFontify( m, ColorBW ) //body // "Func"
  case 'G' : p.CustomFontify( m, ColorBW ) //body
  case 'H' : p.CustomFontify( m, ColorBW ) //body
  case 'I' : p.CustomFontify( m, ColorBW ) //body
  case 'J' : p.CustomFontify( m, ColorBW ) //body
  case 'K' : p.CustomFontify( m, ColorBW ) //body // "keyword"
  case 'L' : p.CustomFontify( m, ColorBW ) //body // "label"
  case 'M' : p.CustomFontify( m, ColorBW ) //body
  case 'N' : p.CustomFontify( m, ColorBW ) //`<span class="defnote" id="` + ToLink(custom) + `" >` + body + "</span>"
  case 'O' : p.CustomFontify( m, ColorBW ) //body
  case 'P' : p.CustomFontify( m, ColorBW ) //body
  case 'Q' : p.CustomFontify( m, ColorBW ) //body
  case 'R' : p.CustomFontify( m, ColorBW ) //body // "result"
  case 'S' : p.CustomFontify( m, ColorBW ) //body
  case 'T' : p.CustomFontify( m, ColorBW ) //body // "radiotarget"
  case 'U' : p.CustomFontify( m, ColorBW ) //body
  case 'V' : p.CustomFontify( m, ColorBW ) //body // "var"
  case 'W' : p.CustomFontify( m, ColorBW ) //body // "warning"
  case 'X' : p.CustomFontify( m, ColorBW ) //body
  case 'Y' : p.CustomFontify( m, ColorBW ) //body
  case 'Z' : p.CustomFontify( m, ColorBW ) //body
  case '\\': p.CustomFontify( m, ColorBW ) //body
  case '^' : p.CustomFontify( m, ColorBW ) //"<sup>" + body + "</sup>"
  case '_' : p.CustomFontify( m, ColorBW ) //"<sub>" + body + "</sub>"
  case '`' : p.CustomFontify( m, ColorBW ) //body
  case 'a' : p.CustomFontify( m, ColorBW ) //"<abbr>" + body + "</abbr>"
  case 'b' : p.CustomFontify( m, ColorBW ) //"<b>" + body + "</b>"
  case 'c' : p.CustomFontify( m, ColorBW ) //"<code>" + body + "</code>"
  case 'd' : p.CustomFontify( m, ColorBW ) //body // data
  case 'e' : p.CustomFontify( m, ColorBW ) //"<em>" + body + "</em>"
  case 'f' : p.CustomFontify( m, ColorBW ) //`<span class="file" >` + body + "</span>"
  case 'g' : p.CustomFontify( m, ColorBW ) //body
  case 'h' : p.CustomFontify( m, ColorBW ) //body
  case 'i' : p.CustomFontify( m, ColorBW ) //"<i>" + body + "</i>"
  case 'j' : p.CustomFontify( m, ColorBW ) //body
  case 'k' : p.CustomFontify( m, ColorBW ) //"<kbd>" + body + "</kbd>"
  case 'l' : p.CustomFontify( m, ColorBW )
//    if custom != "" && custom[ColorBW] == '#' && body != "" && body[ColorBW] == '#' { body = body[1:] }
//    return body//`<a href="` + ToLink( custom ) + `" >` + body + "</a>"
  case 'm' : p.CustomFontify( m, ColorBW ) //`<span class="math" >` + body + "</span>"
  case 'n' : p.CustomFontify( m, ColorBW ) //`<span class="note" ><sup><a href="#` + ToLink(custom) + `" >` + body + "</a></sup></span>"
  case 'o' : p.CustomFontify( m, ColorBW ) //body
  case 'p' : p.CustomFontify( m, ColorBW ) //body
  case 'q' : p.CustomFontify( m, ColorBW )
    // buf.Write( ToU64( "“", '"' ) )
    // buf.Write( body )
    // buf.Write( ToU64( "”", '"' ) )
    // return buf.Data()
  case 'r' : p.CustomFontify( m, ColorBW ) //body // ref
  case 's' : p.CustomFontify( m, ColorBW ) //body
  case 't' : p.CustomFontify( m, ColorBW ) //`<span id="` + custom + `" >` + body + "</span>"
  case 'u' : p.CustomFontify( m, ColorBW ) //"<u>" + body + "</u>"
  case 'v' : p.CustomFontify( m, ColorBW ) //`<code class="verbatim" >` + body + "</code>"
  case 'w' : p.CustomFontify( m, ColorBW ) //body
  case 'x' : p.CustomFontify( m, ColorBW ) //body
  case 'y' : p.CustomFontify( m, ColorBW ) //body
  case 'z' : p.CustomFontify( m, ColorBW ) //body
  case '|' : p.CustomFontify( m, ColorBW ) //body
  case '~' : p.CustomFontify( m, ColorBW ) //body
  }
}

func openDecorator( t uint8 ) []Cell {
  switch t {
  case '"', 'q' : return StrToCells( "“", 0, 0 )
  case '\'': return StrToCells( "‘", 0, 0 )
  case '-' : return StrToCells( "––", 0, 0 )
  }

  return []Cell{}
}

func closeDecorator( t uint8 ) []Cell {
  switch t {
  case '"', 'q' : return StrToCells( "”", 0, 0 )
  case '\'': return StrToCells( "’", 0, 0 )
  case '-' : return StrToCells( "––", 0, 0 )
  }

  return []Cell{}
}

func StrToCells( str string, color, attrs uint8  ) []Cell {
  c := make([]Cell, len( str ) )
  i := 0

  for _, ch := range str {
    c[i] = Cell{ Ch: ch, Color: color, Attrs: attrs }
    i++
  }

  return c[:i]
}
