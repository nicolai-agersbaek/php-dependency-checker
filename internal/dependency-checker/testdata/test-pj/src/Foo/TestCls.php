<?php
declare(strict_types = 1);


namespace TestPrj\Foo;

use \Woo\Something;

/**
 * TODO: Missing class description.
 *
 * @author Nicolai Agersbæk <na@zitcom.dk>
 *
 * @api
 */
class TestCls
{
    
    /**
     * @var Something
     */
    private $err;
    
    public function __construct()
    {
        $this->err = new Something();
    }
}
