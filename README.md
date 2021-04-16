# iso8583-golang
ISO8583 Message Packer and Unpacker with GO

通过配置文件对8583报文进行组包,解包

配置中间中定义了text,number,binary,track，这些分别在yml文件中的type对应。

我们可以通过Default(配置文件路径) 作为全局配置使用。
当我们程序中有多种结构需要使用时（如同时需要64域报文和128域报文）,我们可以通过
NewConfig(配置文件路径)获取配置变量，通过其Pack和UnPack方法实现。

配置文件： 
bit_len为位图长度仅能填写64以及128

fields为各个域属性配置 NORMAL为原值传递,BCDR为BCD右靠,BCDL为BCD左靠,
len_width为长度域长度,如果type为填写则按text处理。
