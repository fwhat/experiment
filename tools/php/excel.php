<?php
class myExcel
{
    /**
     * 获取xlsx 文件的总行数 (因为xlsx转zip解压后,在/xl/worksheets/目录下,是工作表的文本数据,其中含有工作表横纵总数数据)
     * @param $path realpath
     * @return int
     */
    public static function getXlsxCount($path)
    {
        $count = 0;//总数
        $time = time();
        $tempFile = dirname($path) . "/$time.zip";
        copy($path, $tempFile);//将excle复制一份为zip
        $zip = new \ZipArchive;
        $zip->open($tempFile);
        $toDir = $tempFile . $time;
        $zip->extractTo($toDir);//解压
        $xmls = self::getFiles($toDir . '/xl/worksheets/');//获取解压后 /xl/worksheets 路径下的所有文件

        foreach ($xmls as $xml) {
            if (is_dir($xml)) continue;
            $fp = fopen($xml, 'r');
            while (($content = fread($fp, 1024)) !== false) {
                preg_match('/dimension((?!\/).)*/', $content, $res);
                if (isset($res[0])) {
                    if (strrchr($res[0], ':') === false) {
                        $count += 0;
                        break;
                    } else {
                        preg_match('/\d+/', strrchr($res[0], ':'), $res);
                        $count += $res[0];
                        break;
                    }
                }
            }
            fclose($fp);
        }
        unlink($tempFile);//删除压缩文件
        self::delDir($toDir);//删除解压后的文件夹

        return $count;
    }

    /**
     * 获取目录下第一层的文件和目录
     * @param $path
     * @return bool
     */
    public static function getFiles($path)
    {
        $files = [];
        foreach (scandir($path) as $name) {
            if ($name === '.' || $name === '..') {
                continue;
            }
            $files[$name] = rtrim($path , '/') . '/' . $name;
        }

        return $files;
    }

    /**
     * 删除目录
     * @param $dir
     * @return bool
     */
    public static function delDir($dir)
    {
        if (! file_exists($dir)) {
            return true;
        }
        $dir = rtrim($dir, "/") . '/';
        $dirs = scandir($dir);
        foreach ($dirs as $child) {
            if ($child == '.' || $child == '..') continue;
            if (is_dir($dir . $child)) {
                self::delDir($dir . $child);

            } else {
                unlink($dir . $child);
            }
        }
        return rmdir($dir);
    }
}
