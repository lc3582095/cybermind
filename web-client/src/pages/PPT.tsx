import React from 'react';
import { Input, Button, InputNumber, message, Progress, Upload, Tabs, Radio } from 'antd';
import {
  DeleteOutlined,
  FileTextOutlined,
  InfoCircleOutlined,
  UploadOutlined,
  EditOutlined,
  EyeOutlined,
  DownloadOutlined,
  FileMarkdownOutlined,
  FileWordOutlined,
  FilePdfOutlined,
} from '@ant-design/icons';

const { TabPane } = Tabs;
const { TextArea } = Input;

interface PPTTemplate {
  id: string;
  name: string;
  thumbnail: string;
  category: string;
  slides: number;
}

interface PPTGeneration {
  id: number;
  title: string;
  content: string;
  timestamp: Date;
  status: 'generating' | 'completed' | 'error';
  progress: number;
  template: string;
  slideCount: number;
  previewUrl?: string;
}

const TEMPLATE_CATEGORIES = ['全部', '商务', '教育', '科技', '创意', '汇报', '营销'];

const SAMPLE_TEMPLATES: PPTTemplate[] = [
  { id: 't1', name: '商务简约汇报', thumbnail: '/templates/t1.jpg', category: '商务', slides: 12 },
  { id: 't2', name: '科技产品发布', thumbnail: '/templates/t2.jpg', category: '科技', slides: 15 },
  { id: 't3', name: '教育课程分享', thumbnail: '/templates/t3.jpg', category: '教育', slides: 10 },
];

const PPT: React.FC = () => {
  const [activeTab, setActiveTab] = React.useState('template');
  const [title, setTitle] = React.useState('');
  const [content, setContent] = React.useState('');
  const [isGenerating, setIsGenerating] = React.useState(false);
  const [generations, setGenerations] = React.useState<PPTGeneration[]>([]);
  const [selectedTemplate, setSelectedTemplate] = React.useState<string>('');
  const [templateCategory, setTemplateCategory] = React.useState('全部');
  const [slideCount, setSlideCount] = React.useState(10);
  const [importedFile, setImportedFile] = React.useState<any>(null);
  const [inputType, setInputType] = React.useState<'upload' | 'text'>('upload');
  const [inputText, setInputText] = React.useState('');

  const handleGenerate = () => {
    if ((!title.trim() && !importedFile && !inputText.trim()) || isGenerating) return;

    const newGeneration: PPTGeneration = {
      id: Date.now(),
      title: title || importedFile?.name || '未命名PPT',
      content: inputText || '',
      timestamp: new Date(),
      status: 'generating',
      progress: 0,
      template: selectedTemplate,
      slideCount,
    };

    setGenerations([newGeneration, ...generations]);
    setTitle('');
    setContent('');
    setIsGenerating(true);
    setImportedFile(null);
    setInputText('');

    // 模拟生成进度
    let progress = 0;
    const interval = setInterval(() => {
      progress += 10;
      setGenerations((prev) =>
        prev.map((gen) =>
          gen.id === newGeneration.id
            ? {
                ...gen,
                progress,
              }
            : gen
        )
      );

      if (progress >= 100) {
        clearInterval(interval);
        setGenerations((prev) =>
          prev.map((gen) =>
            gen.id === newGeneration.id
              ? {
                  ...gen,
                  status: 'completed',
                  previewUrl: '/preview.jpg',
                }
              : gen
          )
        );
        setIsGenerating(false);
      }
    }, 1000);
  };

  const handleClear = () => {
    setGenerations([]);
    message.success('已清空生成记录');
  };

  const handleFileImport = (info: any) => {
    if (info.file.status === 'done') {
      setImportedFile(info.file);
      message.success(`${info.file.name} 导入成功`);
    }
  };

  const handleInputTypeChange = (e: any) => {
    setInputType(e.target.value);
    setImportedFile(null);
    setInputText('');
  };

  const handleTextInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputText(e.target.value);
  };

  return (
    <div className="flex h-screen bg-[#F7F9FB]">
      {/* 左侧导航栏 */}
      <div className="w-[280px] bg-white border-r border-gray-100 flex flex-col h-screen">
        {/* 顶部标题 */}
        <div className="h-14 px-4 flex items-center border-b border-gray-100 flex-shrink-0">
          <span className="text-base font-medium">AI PPT</span>
          <span className="text-xs text-gray-400 ml-2">@pptgen</span>
        </div>

        {/* 中间可滚动区域 */}
        <div className="flex-1 overflow-y-auto">
          <Tabs activeKey={activeTab} onChange={setActiveTab} className="px-4 pt-4">
            <TabPane tab="模板创建" key="template">
              <div className="space-y-4">
                {/* 模板分类 */}
                <div className="flex flex-wrap gap-2">
                  {TEMPLATE_CATEGORIES.map((category) => (
                    <span
                      key={category}
                      onClick={() => setTemplateCategory(category)}
                      className={`px-3 py-1 rounded-md text-xs cursor-pointer transition-colors ${
                        templateCategory === category
                          ? 'bg-[#EEF2FF] text-[#6366F1]'
                          : 'bg-[#F7F9FB] text-gray-500 hover:bg-gray-100'
                      }`}
                    >
                      {category}
                    </span>
                  ))}
                </div>

                {/* 模板列表 */}
                <div className="grid grid-cols-2 gap-3">
                  {SAMPLE_TEMPLATES
                    .filter((t) => templateCategory === '全部' || t.category === templateCategory)
                    .map((template) => (
                      <div
                        key={template.id}
                        onClick={() => setSelectedTemplate(template.id)}
                        className={`relative rounded-lg overflow-hidden cursor-pointer group ${
                          selectedTemplate === template.id ? 'ring-2 ring-[#6366F1]' : ''
                        }`}
                      >
                        <img
                          src={template.thumbnail}
                          alt={template.name}
                          className="w-full aspect-video object-cover"
                        />
                        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity">
                          <div className="absolute bottom-0 left-0 right-0 p-2 text-white text-xs">
                            <div className="font-medium">{template.name}</div>
                            <div className="text-gray-300">{template.slides}页</div>
                          </div>
                        </div>
                      </div>
                    ))}
                </div>

                {/* PPT标题 */}
                <div>
                  <div className="mb-2">
                    <span className="text-sm">PPT标题</span>
                  </div>
                  <Input
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="请输入PPT标题"
                    className="bg-[#F7F9FB] border-none hover:border-none focus:shadow-none"
                    maxLength={50}
                  />
                </div>

                {/* 幻灯片数量 */}
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm">幻灯片数量</span>
                    <InputNumber
                      value={slideCount}
                      onChange={(value) => setSlideCount(value || 10)}
                      min={5}
                      max={50}
                      className="w-20 bg-[#F7F9FB] border-none"
                      size="small"
                      controls={false}
                    />
                  </div>
                </div>
              </div>
            </TabPane>

            <TabPane tab="文档导入" key="import">
              <div className="space-y-4">
                <Radio.Group value={inputType} onChange={handleInputTypeChange} className="mb-4">
                  <Radio.Button value="upload">文件上传</Radio.Button>
                  <Radio.Button value="text">直接输入</Radio.Button>
                </Radio.Group>

                {inputType === 'upload' ? (
                  <>
                    {/* 文件上传区域 */}
                    <Upload.Dragger
                      accept=".md,.doc,.docx,.txt,.pdf"
                      maxCount={1}
                      onChange={handleFileImport}
                      className="bg-[#F7F9FB] border-dashed border-gray-200 hover:border-[#6366F1]"
                    >
                      <p className="text-gray-500">
                        <UploadOutlined className="text-2xl mb-2" />
                        <br />
                        点击或拖拽上传文档
                      </p>
                      <p className="text-xs text-gray-400 mt-2">
                        支持 Markdown、Word、PDF、文本文件
                      </p>
                    </Upload.Dragger>

                    {/* 支持的格式说明 */}
                    <div className="space-y-2">
                      <div className="flex items-center space-x-2 text-sm text-gray-500">
                        <FileMarkdownOutlined /> <span>Markdown (.md)</span>
                      </div>
                      <div className="flex items-center space-x-2 text-sm text-gray-500">
                        <FileWordOutlined /> <span>Word (.doc, .docx)</span>
                      </div>
                      <div className="flex items-center space-x-2 text-sm text-gray-500">
                        <FilePdfOutlined /> <span>PDF (.pdf)</span>
                      </div>
                      <div className="flex items-center space-x-2 text-sm text-gray-500">
                        <FileTextOutlined /> <span>文本文件 (.txt)</span>
                      </div>
                    </div>

                    {importedFile && (
                      <div className="bg-[#F7F9FB] p-3 rounded-lg">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-2">
                            <FileTextOutlined className="text-[#6366F1]" />
                            <span className="text-sm">{importedFile.name}</span>
                          </div>
                          <Button
                            type="text"
                            icon={<DeleteOutlined />}
                            onClick={() => setImportedFile(null)}
                            className="text-gray-400 hover:text-red-500"
                          />
                        </div>
                      </div>
                    )}
                  </>
                ) : (
                  <div className="space-y-4">
                    <div className="mb-2">
                      <span className="text-sm">输入PPT内容</span>
                      <span className="text-xs text-gray-400 ml-2">支持Markdown格式</span>
                    </div>
                    <TextArea
                      value={inputText}
                      onChange={handleTextInput}
                      placeholder="请输入PPT内容，支持Markdown格式
例如：
# 第一页标题
- 要点1
- 要点2

# 第二页标题
1. 内容1
2. 内容2"
                      className="bg-[#F7F9FB] border-none hover:border-none focus:shadow-none min-h-[300px]"
                      autoSize={{ minRows: 10, maxRows: 20 }}
                    />
                    <div className="text-xs text-gray-400">
                      提示：使用 # 开始新的一页，使用 Markdown 语法添加格式
                    </div>
                  </div>
                )}
              </div>
            </TabPane>
          </Tabs>
        </div>

        {/* 底部生成按钮 */}
        <div className="flex-shrink-0 p-4 border-t border-gray-100 bg-white">
          <Button
            type="primary"
            block
            onClick={handleGenerate}
            loading={isGenerating}
            className="bg-[#6366F1] hover:bg-[#5558E6]"
            disabled={(!title && !importedFile && !inputText.trim()) || isGenerating}
          >
            {isGenerating ? '生成中...' : '生成'}
          </Button>
          <div className="text-xs text-gray-400 text-center mt-2">
            {activeTab === 'template' ? '选择模板后开始生成' : (inputType === 'upload' ? '导入文档后开始生成' : '输入内容后开始生成')}
          </div>
        </div>
      </div>

      {/* 右侧内容区域 */}
      <div className="flex-1 flex flex-col">
        {/* 顶部标题栏 */}
        <div className="h-14 px-6 bg-white border-b border-gray-100 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium">AI PPT</span>
            <span className="text-xs text-gray-400">@pptgen</span>
          </div>
          <div className="flex items-center space-x-2">
            <Button size="small" type="text" icon={<InfoCircleOutlined />}>
              参考教程
            </Button>
          </div>
        </div>

        {/* 主要内容区域 */}
        <div className="flex-1 overflow-y-auto p-6">
          {generations.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-gray-400">
              <FileTextOutlined className="text-4xl mb-4" />
              <p>暂无PPT创作</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-6">
              {generations.map((generation) => (
                <div
                  key={generation.id}
                  className="bg-white rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow"
                >
                  {/* 预览图 */}
                  <div className="relative aspect-video bg-gray-50">
                    {generation.status === 'generating' ? (
                      <div className="absolute inset-0 flex flex-col items-center justify-center">
                        <FileTextOutlined className="text-2xl text-[#6366F1] mb-2" />
                        <Progress
                          type="circle"
                          percent={generation.progress}
                          width={60}
                          strokeWidth={4}
                          strokeColor={{ '0%': '#818cf8', '100%': '#6366f1' }}
                        />
                      </div>
                    ) : (
                      <>
                        <img
                          src={generation.previewUrl}
                          alt={generation.title}
                          className="w-full h-full object-cover"
                        />
                        <div className="absolute inset-0 bg-black/30 opacity-0 group-hover:opacity-100 transition-opacity">
                          <div className="absolute inset-0 flex items-center justify-center space-x-2">
                            <Button
                              type="text"
                              icon={<EyeOutlined />}
                              className="text-white hover:text-[#6366F1]"
                            >
                              预览
                            </Button>
                            <Button
                              type="text"
                              icon={<EditOutlined />}
                              className="text-white hover:text-[#6366F1]"
                            >
                              编辑
                            </Button>
                          </div>
                        </div>
                      </>
                    )}
                  </div>

                  {/* 信息区域 */}
                  <div className="p-4">
                    <div className="flex items-center justify-between mb-2">
                      <h3 className="text-sm font-medium">{generation.title}</h3>
                      <div className="flex items-center space-x-2">
                        <Button
                          type="text"
                          icon={<DownloadOutlined />}
                          className="text-gray-400 hover:text-[#6366F1]"
                        />
                        <Button
                          type="text"
                          icon={<DeleteOutlined />}
                          className="text-gray-400 hover:text-red-500"
                        />
                      </div>
                    </div>
                    <div className="flex items-center space-x-2 text-xs text-gray-400">
                      <span>{generation.slideCount}页</span>
                      <span>·</span>
                      <span>{generation.timestamp.toLocaleTimeString()}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default PPT; 