import React from 'react';
import { Card, Button, Tag, Input } from 'antd';
import {
  SearchOutlined,
  AppstoreOutlined,
  RobotOutlined,
  PictureOutlined,
  VideoCameraOutlined,
  CustomerServiceOutlined,
  FileTextOutlined,
  TranslationOutlined,
  CodeOutlined,
  EditOutlined,
  CalculatorOutlined,
  BulbOutlined,
  StarOutlined,
} from '@ant-design/icons';

const { Search } = Input;

interface AppItem {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  tags: string[];
  isNew?: boolean;
  isHot?: boolean;
}

const apps: AppItem[] = [
  {
    id: 'chat',
    name: 'AI对话',
    description: '智能对话助手，可以回答问题、写文章、做分析等',
    icon: <RobotOutlined className="text-2xl" />,
    tags: ['对话', '写作', '助手'],
    isHot: true,
  },
  {
    id: 'image',
    name: 'AI绘画',
    description: '智能图像生成，可以创作各种风格的图片',
    icon: <PictureOutlined className="text-2xl" />,
    tags: ['绘画', '设计', '创作'],
    isHot: true,
  },
  {
    id: 'video',
    name: 'AI视频',
    description: '智能���频生成，可以制作短视频、动画等',
    icon: <VideoCameraOutlined className="text-2xl" />,
    tags: ['视频', '动画', '创作'],
    isNew: true,
  },
  {
    id: 'music',
    name: 'AI音乐',
    description: '智能音乐生成，可以创作音乐、音效等',
    icon: <CustomerServiceOutlined className="text-2xl" />,
    tags: ['音乐', '音效', '创作'],
    isNew: true,
  },
  {
    id: 'ppt',
    name: 'AI PPT',
    description: '智能PPT生成，快速制作精美演示文稿',
    icon: <FileTextOutlined className="text-2xl" />,
    tags: ['PPT', '设计', '办公'],
    isHot: true,
  },
  {
    id: 'translate',
    name: 'AI翻译',
    description: '智能多语言翻译，支持多种语言互译',
    icon: <TranslationOutlined className="text-2xl" />,
    tags: ['翻译', '语言', '工具'],
  },
  {
    id: 'code',
    name: 'AI编程',
    description: '智能代码生成，帮助编写和优化代码',
    icon: <CodeOutlined className="text-2xl" />,
    tags: ['编程', '开发', '工具'],
  },
  {
    id: 'write',
    name: 'AI写作',
    description: '智能写作助手，帮助创作各类文章',
    icon: <EditOutlined className="text-2xl" />,
    tags: ['写作', '创作', '助手'],
  },
  {
    id: 'math',
    name: 'AI数学',
    description: '智能数学解题，帮助解决数学问题',
    icon: <CalculatorOutlined className="text-2xl" />,
    tags: ['数学', '教育', '工具'],
  },
  {
    id: 'idea',
    name: 'AI创意',
    description: '智能创意生成，激发创新思维',
    icon: <BulbOutlined className="text-2xl" />,
    tags: ['创意', '创新', '助手'],
  },
];

const Apps: React.FC = () => {
  return (
    <div className="w-[calc(100vw-320px)] fixed right-0 top-0 bottom-0 bg-gray-50">
      {/* 顶部标题栏 */}
      <div className="fixed top-0 right-0 w-[calc(100vw-320px)] h-14 px-4 border-b border-gray-100 bg-white flex items-center justify-between z-10">
        <div className="flex items-center space-x-2">
          <span className="text-base font-medium">AI应用</span>
        </div>
        <div className="flex items-center space-x-2">
          <Search
            placeholder="搜索应用"
            style={{ width: 200 }}
            className="search-input"
          />
        </div>
      </div>

      {/* 应用列表 */}
      <div className="fixed top-14 right-0 bottom-0 w-[calc(100vw-320px)] overflow-y-auto">
        <div className="max-w-6xl mx-auto p-4">
          <div className="grid grid-cols-3 gap-4">
            {apps.map((app) => (
              <Card
                key={app.id}
                className="group hover:shadow-md transition-shadow"
                bodyStyle={{ padding: '1.5rem' }}
              >
                <div className="flex items-start space-x-4">
                  <div className="w-12 h-12 rounded-lg bg-primary-50 flex items-center justify-center text-primary-500">
                    {app.icon}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-2">
                      <h3 className="text-base font-medium text-gray-900">{app.name}</h3>
                      {app.isNew && (
                        <Tag color="success" className="rounded-full">
                          New
                        </Tag>
                      )}
                      {app.isHot && (
                        <Tag color="error" className="rounded-full">
                          Hot
                        </Tag>
                      )}
                    </div>
                    <p className="mt-1 text-sm text-gray-500 line-clamp-2">
                      {app.description}
                    </p>
                    <div className="mt-3 flex items-center flex-wrap gap-2">
                      {app.tags.map((tag) => (
                        <Tag key={tag} className="rounded-full">
                          {tag}
                        </Tag>
                      ))}
                    </div>
                  </div>
                </div>
                <div className="mt-4 flex items-center justify-end">
                  <Button type="primary" className="bg-primary-500">
                    立即使用
                  </Button>
                </div>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Apps; 