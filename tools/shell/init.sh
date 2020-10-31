#!/bin/bash
start_time=`date +%s`
function _echo()
{
	echo -e "\033[40;35m $1 \033[0m"
}
_echo "请输入要创建的常用用户名(平时使用应避免使用root):"
read u_name
_echo "是否禁止root远程登录(建议禁止) y/n"
read allow_root
_echo ">>>>脚本开始执行，无需再做操作，估计耗时5分钟，请耐心等待..."
sleep 2
_echo ">>>>正在配置防火墙"
c_centos=`cat /etc/redhat-release | sed 's/[a-zA-Z ]*\([0-9]\)\([0-9a-zA-Z\. ()]\)\+/\1/g'`
c_ssh_port=$(($RANDOM%4000+4000))
case ${c_centos} in
    6)  iptables -I INPUT -p tcp --dport 80 -j ACCEPT
		iptables -I INPUT -p tcp --dport 22 -j ACCEPT
		iptables -I INPUT -p tcp --dport 443 -j ACCEPT
		iptables -I INPUT -p tcp --dport 3306 -j ACCEPT
		iptables -I INPUT -p tcp --dport ${c_ssh_port} -j ACCEPT
    ;;
    7)  systemctl start firewalld.service
		firewall-cmd --zone=public --add-port=80/tcp --permanent > /dev/null 2>&1
		firewall-cmd --zone=public --add-port=22/tcp --permanent > /dev/null 2>&1
		firewall-cmd --zone=public --add-port=3306/tcp --permanent > /dev/null 2>&1
		firewall-cmd --zone=public --add-port=443/tcp --permanent > /dev/null 2>&1
		firewall-cmd --zone=public --add-port=${c_ssh_port}/tcp --permanent > /dev/null 2>&1
		firewall-cmd --reload
		systemctl enable firewalld.service
    ;;
    *)  _echo "系统版本太低或不是centos系统"
		exit 1
    ;;
esac
_echo ">>>>正在更改yum源至阿里云源"
sleep 1
mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
case ${c_centos} in
    6)  curl -o -s /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-6.repo
    ;;
    7)  curl -o -s /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
    ;;
    *)  _echo "系统版本太低或不是centos系统"
		exit 1
    ;;
esac
yum clean all
yum -y makecache
_echo ">>>>正在更新yum包库..."
sleep 1
rpm -Uh https://mirror.webtatic.com/yum/el7/epel-release.rpm
rpm -Uh https://mirror.webtatic.com/yum/el7/webtatic-release.rpm
yum -y update
_echo ">>>>正在检测脚本需要的依赖..."
sleep 1
command -v vim > /dev/null 2>&1 || c_arr[${#c_arr[*]}]="vim"
command -v git > /dev/null 2>&1 || c_arr[${#c_arr[*]}]="git"
command -v wget > /dev/null 2>&1 || c_arr[${#c_arr[*]}]="wget"
_echo ">>>>需要安装的依赖有:${c_arr}"
_echo ">>>>执行安装..."
sleep 1
for i in ${c_arr[*]}
do
	yum -y install ${i}
done
_echo ">>>>正在创建常用用户..."
useradd -g users -m ${u_name}
u_passwd=`echo $(date)${u_name}|base64`
u_passwd=${u_passwd:${#u_passwd}-16:${#u_passwd}}
echo ${u_passwd} | passwd ${u_name} --stdin > /dev/null
echo "${u_name} ALL=(ALL) ALL" >> /etc/sudoers
sed -i 's/#Port/Port/g' /etc/ssh/sshd_config
sed -i "s/Port 22/Port ${c_ssh_port}/g" /etc/ssh/sshd_config
if [ ${allow_root}='y' ] || [ ${allow_root}='Y' ]; then
	sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/g' /etc/ssh/sshd_config
	sed -i 's/PermitRootLogin yes/PermitRootLogin no/g' /etc/ssh/sshd_config
fi
service sshd restart
_echo ">>>>正在安装zsh..."
sleep 1
yum -y install zsh
chsh -s /bin/zsh ${u_name} 
_echo ">>>>正在安装oh-my-zsh..."
sleep 1
su - ${u_name} -c "wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | sh"
sed -i 's/ZSH_THEME=\"robbyrussell\"/ZSH_THEME=\"ys\"/g' /home/${u_name}/.zshrc
_echo ">>>>正在安装php71..."
sleep 1
yum -y install mod_php71w php71w-bcmath php71w-cli php71w-common php71w-devel php71w-fpm php71w-gd php71w-mbstring php71w-mcrypt php71w-mysql php71w-snmp  php71w-xml php71w-process php71w-ldap net-snmp net-snmp-devel net-snmp-utils rrdtool
_echo ">>>>正在安装nginx..."
sleep 1
yum -y install nginx
echo "<?php phpinfo();" > /usr/share/nginx/html/test.php
cat >> /etc/nginx/default.d/test.conf <<EOF
location ~ .*\.php.* {
        fastcgi_pass  127.0.0.1:9000;
        fastcgi_index   index.php;
        include fastcgi.conf;
}
EOF
nginx
php-fpm
u_ip=`curl -s ident.me`
_echo ">>>>脚本执行完毕"
_echo "成功安装:${c_arr} nginx php zsh"
_echo "访问${u_ip}:/test.php 测试nginx和php是否安装正常"
_echo "注: 测试后请及时删除测试文件/usr/share/nginx/html/test.php和/etc/nginx/default.d/test.conf"
_echo "注: 如服务器有安全组,请自行在安全组配置对应端口,否则可能无法正常访问站点(如阿里云,aws...)"
_echo "常用用户名:${u_name} 密码:${u_passwd}"
_echo "请牢记密码, 或修改密码(passwd ${u_name})"
_echo "防火墙开放端口: 22, 80, 443, 3306, ${c_ssh_port}(新的ssh端口)"
if [ ${allow_root}='y' ] || [ ${allow_root}='Y' ]; then
	_echo "提示: 已禁止root远程登录"
	_echo "请先用新创建的用户再其他终端成功连接后，再关闭本连接"
fi
_echo "脚本执行耗时$(expr `date +%s` - ${start_time})s"





