<?php
/**
 * Created by PhpStorm.
 * User: Mr.Zhou
 * Date: 2018/11/12
 * Time: 2:09 PM
 */
defined('APP_DIR') || define('APP_DIR', dirname(__FILE__));

require_once APP_DIR . '/vendor/autoload.php';

try {
    $client = new \Sms\SmsClient("127.0.0.1:50051", [
        'credentials' => Grpc\ChannelCredentials::createInsecure()
    ]);
    $request = new \Sms\Request();
    $request->setParams(json_encode(['code' => '123456']));
    $request->setMobile(15600087538);
    //$request->setTemplate("123");
    list($reply, $error) = $client->Send($request)->wait();
    if ($error->code) {
        exit("发送失败，错误信息：" . $error->details);
    }
    var_dump($reply->getCode(), $reply->getMsg(), $reply->serializeToJsonString());

} catch (\Exception $exception) {
    var_dump($exception->getMessage());
}