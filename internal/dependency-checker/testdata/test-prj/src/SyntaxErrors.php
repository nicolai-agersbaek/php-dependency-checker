<?php
declare(strict_types = 1);


namespace TestPrj;

use TestPrj\SomeBar\SomeCls;
use TestPrj\Somebar\SomeCls as SomeClsAlias;

/**
 * Class for testing references to namespaces using incorrect case.
 *
 * @author Nicolai Agersbæk <na@zitcom.dk>
 *
 * @api
 */
class SyntaxErrors
{
    
    /**
     * @param SomeClsAlias $some
     * @param SomeCls      $someCls
     */
    public function __construct(SomeClsAlias $some, $someCls SomeCls)
    {
        $a
    }
}
}
