<template>
    <header id="blog-header">
        <div class="hdr-wrapper section-inner">
            <div class="hdr-left">
                <div class="site-branding">
                    <a href="#">Hola-CosMos</a>
                </div>
            </div>
            <div class="hdr-right">
                <Dropdown>
                    <a href="javascript:void(0)">
                        下拉菜单
                    </a>
                    <DropdownMenu slot="list">
                        <DropdownItem>驴打滚</DropdownItem>
                        <DropdownItem>炸酱面</DropdownItem>
                        <DropdownItem disabled>豆汁儿</DropdownItem>
                        <DropdownItem>冰糖葫芦</DropdownItem>
                        <DropdownItem divided>北京烤鸭</DropdownItem>
                    </DropdownMenu>
                </Dropdown>
                <Divider type="vertical" />
                <Dropdown>
                    <a href="javascript:void(0)">
                        theme
                        <Icon type="ios-arrow-down"></Icon>
                    </a>
                    <DropdownMenu slot="list">
                        <DropdownItem>驴打滚</DropdownItem>
                        <DropdownItem>炸酱面</DropdownItem>
                        <DropdownItem disabled>豆汁儿</DropdownItem>
                        <DropdownItem>冰糖葫芦</DropdownItem>
                        <DropdownItem divided>北京烤鸭</DropdownItem>
                    </DropdownMenu>
                </Dropdown>
            </div>
        </div>
    </header>
</template>

<script>
export default {

}
</script>
