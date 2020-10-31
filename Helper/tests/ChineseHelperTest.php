<?php
namespace Tests;

class ChineseHelperTest extends \PHPUnit\Framework\TestCase
{
    public function testGetFirstLetter()
    {
        $helper = new ChineseHelper();
        $helper->setChinese('中文');

        $this->assertTrue('Z' === $this->testGetFirstLetter());
    }
}


//require __DIR__ . '/../src/ChineseHelper.php';
//
//$helper = new ChineseHelper();
//
//$helper->setChinese('中文');
//echo $helper->getFirstLetter();