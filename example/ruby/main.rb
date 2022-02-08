def fibonacci( n )
    [ n ] if ( 0..1 ).include? n
    ( fibonacci( n - 1 ) + fibonacci( n - 2 ) ) if n > 1
end

puts fibonacci( 5 )
