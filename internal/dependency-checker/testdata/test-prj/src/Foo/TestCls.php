<?php
/** @noinspection PhpUndefinedNamespaceInspection */
/** @noinspection PhpUndefinedClassInspection */
/** @noinspection PhpMissingDocCommentInspection */
declare(strict_types = 1);


namespace TestPrj\Foo;

use Woo\Something;

class TestCls extends \Some\Parent implements \Foo\SomeInterface, Bar\OtherInterface
{
    
    public function __construct()
    {
        $foo = \is_string('foo');
        $bar = is_string('bar');
        
        $something = new Something();
        $else = new \Bar\SomethingElse();
    }
}
