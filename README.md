# Boss爬虫## 招聘- [x] 筛选简历,过滤已被同事撩过的- [x] 对候选人进行打分，排序- [x] 自动打招呼- [x] 自动索要简历- [x] 飞书自动通知## Quick Start 1.编译```bashgo build -o boss *.go```2. 启动```bash./boss ```3. 设置定时任务```bashcrontab -e# 新增30 09 * * * /var/boss > /var/boss.log```### TODO- [ ] 支持读取Mysql配置- [ ] 支持多账号- [ ] 支持多Job同时招聘- [ ] 每个职位打招呼人数可配置- [ ] 支持http api- [ ] 支持内部定时任务