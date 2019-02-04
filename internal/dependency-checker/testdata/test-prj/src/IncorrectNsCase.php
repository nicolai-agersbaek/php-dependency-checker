<?php
declare(strict_types = 1);


namespace TestPrj;

use TestPrj\Somebar\SomeCls;

/**
 * Class for testing references to namespaces using incorrect case.
 *
 * @author Nicolai Agersbæk <na@zitcom.dk>
 *
 * @api
 */
class IncorrectNsCase
{
    
    /**
     * @param SomeCls $some
     */
    public function __construct(SomeCls $some)
    {
    }
}
