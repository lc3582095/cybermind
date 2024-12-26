import React, { useState } from 'react';
import { Input, Button, Select, Upload, message, Slider, Radio, Tooltip, Tabs, Card, Tag, Drawer, List, Divider, Switch, Modal, Space } from 'antd';
import {
  SendOutlined, DeleteOutlined, UploadOutlined, DownloadOutlined,
  PictureOutlined, LoadingOutlined, PlusOutlined, SettingOutlined,
  QuestionCircleOutlined, EditOutlined, CopyOutlined, UndoOutlined,
  HistoryOutlined, StarOutlined, StarFilled, SaveOutlined,
  TranslationOutlined, SwapOutlined, FileImageOutlined,
  ExpandOutlined, CompressOutlined, CheckOutlined, MinusOutlined
} from '@ant-design/icons';
import type { UploadFile } from 'antd/es/upload/interface';

const { TextArea } = Input;
const { TabPane } = Tabs;

interface ImageGeneration {
  id: number;
  prompt: string;
  negativePrompt: string;
  imageUrl: string;
  timestamp: Date;
  status: 'generating' | 'completed' | 'error';
  model: string;
  width: number;
  height: number;
  steps: number;
  seed: number;
  cfg: number;
}

interface ModelConfig {
  name: string;
  label: string;
  description: string;
  features: string[];
  icon?: React.ReactNode;
}

const MODELS: ModelConfig[] = [
  {
    name: 'midjourney',
    label: 'Midjourney',
    description: '真实风格，艺术感强',
    features: ['高质量艺术创作', '真实感渲染', '细节丰富'],
  },
  {
    name: 'niji',
    label: 'Niji Journey',
    description: '动漫风格',
    features: ['动漫风格', '角色设计', '场景绘制'],
  },
  {
    name: 'dalle3',
    label: 'DALL·E 3',
    description: '自然语言理解能力强',
    features: ['精确的文本理解', '创意构图', '多样化输出'],
  },
  {
    name: 'sdxl',
    label: 'Stable Diffusion XL',
    description: '开源模型，可调参数多',
    features: ['高度可定制', '本地部署', '多种风格'],
  },
];

const ASPECT_RATIOS = [
  { label: '1:1', value: '1:1', width: 1024, height: 1024 },
  { label: '3:2', value: '3:2', width: 1024, height: 683 },
  { label: '4:3', value: '4:3', width: 1024, height: 768 },
  { label: '16:9', value: '16:9', width: 1024, height: 576 },
  { label: '9:16', value: '9:16', width: 576, height: 1024 },
];

const STYLE_PRESETS = [
  { label: '写实风格', value: 'realistic' },
  { label: '油画风格', value: 'oil-painting' },
  { label: '水彩风格', value: 'watercolor' },
  { label: '动漫风格', value: 'anime' },
  { label: '素描风格', value: 'sketch' },
  { label: '3D渲染', value: '3d-render' },
];

const PROMPT_TEMPLATES = [
  {
    category: '风格',
    items: [
      { label: '史诗级震撼', value: 'epic, cinematic, dramatic lighting, volumetric lighting' },
      { label: '梦幻童话', value: 'dreamy, fairy tale, soft lighting, magical, whimsical' },
      { label: '赛博朋克', value: 'cyberpunk, neon, futuristic, sci-fi, dark, rainy' },
      { label: '水彩画', value: 'watercolor, artistic, flowing colors, soft edges' },
    ]
  },
  {
    category: '场景',
    items: [
      { label: '自然风光', value: 'landscape, nature, mountains, forest, scenic' },
      { label: '城市街景', value: 'cityscape, urban, street view, architecture' },
      { label: '太空科幻', value: 'space, stars, nebula, planets, cosmic' },
      { label: '海底世界', value: 'underwater, ocean, marine life, coral reef' },
    ]
  },
  {
    category: '光照',
    items: [
      { label: '黄金时刻', value: 'golden hour, warm lighting, sunset' },
      { label: '蓝调时分', value: 'blue hour, cool tones, twilight' },
      { label: '体积光', value: 'volumetric lighting, god rays, atmospheric' },
      { label: '霓虹灯', value: 'neon lighting, colorful glow, night scene' },
    ]
  }
];

const NEGATIVE_TEMPLATES = [
  { label: '常见问题', value: 'blurry, duplicate, deformed, bad anatomy, disfigured, poorly drawn, extra limbs, text, watermark' },
  { label: '画质问题', value: 'low quality, pixelated, artifacts, jpeg artifacts, noise, grain' },
  { label: '构图问题', value: 'bad composition, unbalanced, cropped, frame, border' },
];

const QUICK_PROMPTS_CATEGORIES = [
  {
    label: '质量',
    items: [
      { label: '高清', value: 'high quality, 8k, ultra detailed', count: 0 },
      { label: '精致', value: 'masterpiece, best quality, extremely detailed', count: 0 },
    ]
  },
  {
    label: '风格',
    items: [
      { label: '电影感', value: 'cinematic, dramatic lighting', count: 0 },
      { label: '写实', value: 'photorealistic, hyperrealistic', count: 0 },
      { label: '概念艺术', value: 'concept art, digital art', count: 0 },
      { label: '油画', value: 'oil painting, masterpiece', count: 0 },
    ]
  },
  {
    label: '光照',
    items: [
      { label: '体积光', value: 'volumetric lighting, god rays', count: 0 },
      { label: '柔光', value: 'soft lighting, ambient light', count: 0 },
    ]
  }
];

const Image: React.FC = () => {
  const [prompt, setPrompt] = useState('');
  const [negativePrompt, setNegativePrompt] = useState('');
  const [selectedModel, setSelectedModel] = useState('midjourney');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generations, setGenerations] = useState<ImageGeneration[]>([
    {
      id: 1,
      prompt: "示例图片1",
      negativePrompt: "",
      imageUrl: "https://picsum.photos/400",
      status: "completed",
      model: "Stable Diffusion",
      timestamp: new Date(),
      width: 512,
      height: 512,
      steps: 20,
      seed: -1,
      cfg: 7
    },
    {
      id: 2,
      prompt: "示例图片2",
      negativePrompt: "",
      imageUrl: "https://picsum.photos/401",
      status: "completed",
      model: "Stable Diffusion",
      timestamp: new Date(),
      width: 512,
      height: 512,
      steps: 20,
      seed: -1,
      cfg: 7
    }
  ]);
  const [referenceImage, setReferenceImage] = useState<UploadFile | null>(null);
  const [aspectRatio, setAspectRatio] = useState('1:1');
  const [stylePreset, setStylePreset] = useState('realistic');
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [steps, setSteps] = useState(30);
  const [cfg, setCfg] = useState(7);
  const [seed, setSeed] = useState(-1);
  const [showHistory, setShowHistory] = useState(false);
  const [favorites, setFavorites] = useState<number[]>([]);
  const [showTemplates, setShowTemplates] = useState(false);
  const [autoTranslate, setAutoTranslate] = useState(true);
  const [showFullscreen, setShowFullscreen] = useState(false);
  const [selectedImage, setSelectedImage] = useState<ImageGeneration | null>(null);
  const [showReferencePreview, setShowReferencePreview] = useState(false);
  const [promptHistory, setPromptHistory] = useState<string[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [referenceStrength, setReferenceStrength] = useState(50);
  const [version, setVersion] = useState('v6.1');
  const [quality, setQuality] = useState('ultra');
  const [style, setStyle] = useState('realistic');
  const [viewMode, setViewMode] = useState('grid');
  const [chaos, setChaos] = useState(50);
  const [stylization, setStylization] = useState(50);

  const handleGenerate = () => {
    if (!prompt.trim()) return;

    const selectedRatio = ASPECT_RATIOS.find(ratio => ratio.value === aspectRatio);
    const newGeneration: ImageGeneration = {
      id: Date.now(),
      prompt: prompt.trim(),
      negativePrompt: negativePrompt,
      imageUrl: "",
      status: "generating",
      model: selectedModel,
      timestamp: new Date(),
      width: selectedRatio?.width || 512,
      height: selectedRatio?.height || 512,
      steps: steps,
      seed: seed,
      cfg: cfg
    };
    
    setGenerations(prev => [newGeneration, ...prev]);
    setIsGenerating(true);

    // 模拟生成过程
    setTimeout(() => {
      setGenerations(prev => 
        prev.map(gen => 
          gen.id === newGeneration.id
            ? { ...gen, status: "completed", imageUrl: `https://picsum.photos/${400 + gen.id % 10}` }
            : gen
        )
      );
      setIsGenerating(false);
    }, 2000);
  };

  const handleClear = () => {
    Modal.confirm({
      title: "确认清空历史记录？",
      content: "清空后将无法恢复",
      onOk: () => {
    setGenerations([]);
        message.success("已清空历史记录");
      }
    });
  };

  const handleVariation = (generation: ImageGeneration) => {
    // TODO: 实现图片变体生成
    message.info('正在生成变体...');
  };

  const handleUpscale = (generation: ImageGeneration) => {
    // TODO: 实现图片放大
    message.info('正在放大图片...');
  };

  const handleInpaint = (generation: ImageGeneration) => {
    // TODO: 实现局部重绘
    message.info('请选择要重绘的区域...');
  };

  const handleToggleFavorite = (id: number) => {
    setFavorites(prev => 
      prev.includes(id) ? prev.filter(fid => fid !== id) : [...prev, id]
    );
  };

  const handleApplyTemplate = (template: string) => {
    setPrompt(prev => prev ? `${prev}, ${template}` : template);
  };

  const handleSavePreset = () => {
    if (!prompt) return;
    // TODO: 保存当前设置为预设
    message.success('已保存为预设');
  };

  const handlePromptSubmit = () => {
    if (!prompt.trim()) return;
    setPromptHistory(prev => [prompt, ...prev.slice(0, 9)]);
    handleGenerate();
  };

  return (
    <div className="flex flex-col h-full">
      {/* 顶部标题栏 */}
      <div className="h-14 px-4 border-b border-gray-100 bg-white flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <span className="text-base font-medium">AI绘画</span>
          <Tag color="blue">设计中</Tag>
        </div>
        <div className="flex items-center space-x-2">
          <Button 
            type="text" 
            icon={<QuestionCircleOutlined />}
            className="text-gray-400"
          >
            使用教程
          </Button>
          <Button 
            type={showAdvanced ? 'primary' : 'default'}
            icon={<SettingOutlined />}
            onClick={() => setShowAdvanced(!showAdvanced)}
          >
            高级设置
          </Button>
        </div>
      </div>

      {/* 主要内容区域 */}
      <div className="flex-1 bg-gray-50 p-4 lg:p-6 overflow-auto">
        <div className="flex flex-col lg:flex-row h-full gap-4 lg:gap-6">
          {/* 左侧内容区域 */}
          <div className="w-full lg:w-[360px] flex flex-col gap-4">
            {/* 提示词输入区域 */}
            <div className="bg-white rounded-lg shadow-sm">
              <div className="p-4 border-b border-gray-100">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <Select
                      value={selectedModel}
                      onChange={setSelectedModel}
                      style={{ width: 160 }}
                      options={MODELS.map(model => ({
                        value: model.name,
                        label: (
                          <div className="flex flex-col">
                            <span className="text-sm">{model.label}</span>
                            <span className="text-xs text-gray-400">{model.description}</span>
                          </div>
                        )
                      }))}
                    />
                    <Select
                      value={version}
                      style={{ width: 80 }}
                      options={[
                        { label: 'V6.1', value: 'v6.1' },
                        { label: 'V6', value: 'v6' },
                        { label: 'V5.2', value: 'v5.2' },
                        { label: 'V5.1', value: 'v5.1' }
                      ]}
                    />
                  </div>
                  <Button
                    type={showAdvanced ? 'primary' : 'default'}
                    size="small"
                    icon={<SettingOutlined />}
                    onClick={() => setShowAdvanced(!showAdvanced)}
                  >
                    高级设置
                  </Button>
                </div>
              </div>

              <Tabs defaultActiveKey="prompt" className="text-sm">
                <TabPane tab="正向提示词" key="prompt">
                  <div className="p-4">
                    <div className="relative">
                      <TextArea
                        value={prompt}
                        onChange={(e) => setPrompt(e.target.value)}
                        placeholder="描述你想要生成的图片内容..."
                        autoSize={{ minRows: 3, maxRows: 6 }}
                        className="mb-3 text-sm pr-16"
                        maxLength={500}
                      />
                      <div className="absolute right-2 bottom-5 text-xs text-gray-400">
                        {prompt.length}/500
                      </div>
                    </div>
                    <div className="flex items-center justify-between mb-2">
                      <div className="text-xs text-gray-500">快速添加</div>
                      <div className="flex items-center space-x-2">
                        <Switch
                          size="small"
                          checked={autoTranslate}
                          onChange={setAutoTranslate}
                        />
                        <span className="text-xs text-gray-500">自动翻译</span>
                      </div>
                    </div>
                    <div className="flex flex-wrap gap-1.5">
                      {QUICK_PROMPTS_CATEGORIES.slice(0, 2).reduce((acc: JSX.Element[], category) => [
                        ...acc,
                        ...category.items.map((item: { label: string; value: string; count: number }) => (
                          <Tag
                            key={item.label}
                            className="cursor-pointer text-xs hover:bg-primary-50 transition-colors"
                            onClick={() => handleApplyTemplate(item.value)}
                          >
                            {item.label}
                          </Tag>
                        ))
                      ], [])}
                      <Button
                        type="link"
                        size="small"
                        className="text-xs px-2 py-0"
                        icon={<PlusOutlined className="text-xs" />}
                        onClick={() => setShowTemplates(true)}
                      >
                        更多模板
                      </Button>
                    </div>
                  </div>
                </TabPane>
                <TabPane tab="负向提示词" key="negative">
                  <div className="p-4">
                    <div className="relative">
                      <TextArea
                        value={negativePrompt}
                        onChange={(e) => setNegativePrompt(e.target.value)}
                        placeholder="描述你不想在图片中出现的内容..."
                        autoSize={{ minRows: 3, maxRows: 6 }}
                        className="mb-3 text-sm pr-16"
                        maxLength={500}
                      />
                      <div className="absolute right-2 bottom-5 text-xs text-gray-400">
                        {negativePrompt.length}/500
                      </div>
                    </div>
                    <div className="flex flex-wrap gap-1.5">
                      {NEGATIVE_TEMPLATES.map(item => (
                        <Tag
                          key={item.label}
                          className="cursor-pointer text-xs hover:bg-primary-50 transition-colors"
                          onClick={() => setNegativePrompt(item.value)}
                        >
                          {item.label}
                        </Tag>
                      ))}
                    </div>
                  </div>
                </TabPane>
              </Tabs>
            </div>

            {/* 参数设置区域 */}
            <div className="bg-white rounded-lg shadow-sm">
              <div className="p-4 border-b border-gray-100">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">基础参数</span>
                  <Button
                    type={showAdvanced ? 'primary' : 'default'}
                    size="small"
                    icon={<SettingOutlined />}
                    onClick={() => setShowAdvanced(!showAdvanced)}
                  >
                    高级设置
                  </Button>
                </div>
              </div>

              <div className="p-4">
                <div className="space-y-6">
                  {/* 图片比例设置 */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">图片比例</span>
                        <Tooltip title="选择生成图片的宽高比">
                          <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                        </Tooltip>
                      </div>
                      <div className="text-xs text-gray-400">
                        推荐: 1:1, 16:9
                      </div>
                    </div>
                    <Radio.Group
                      value={aspectRatio}
                      onChange={e => setAspectRatio(e.target.value)}
                      className="w-full grid grid-cols-5 gap-1"
                    >
                      {ASPECT_RATIOS.map(ratio => (
                        <Radio.Button 
                          key={ratio.value} 
                          value={ratio.value} 
                          className="text-xs text-center"
                          style={{ padding: '4px 0' }}
                        >
                          {ratio.label}
                        </Radio.Button>
                      ))}
                    </Radio.Group>
                  </div>

                  {/* 品质与风格设置 */}
                  <div className="grid grid-cols-2 gap-6">
                    <div>
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-sm font-medium">品质</span>
                        <Tooltip title="更高品质需要更多算力">
                          <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                        </Tooltip>
                      </div>
                      <Radio.Group 
                        value={quality} 
                        onChange={e => setQuality(e.target.value)}
                        className="flex flex-col gap-2"
                      >
                        <Radio value="ultra" className="text-sm">
                          <div className="flex items-center justify-between w-full">
                            <span>超高清</span>
                            <span className="text-xs text-gray-400">1024×1024</span>
                          </div>
                        </Radio>
                        <Radio value="high" className="text-sm">
                          <div className="flex items-center justify-between w-full">
                            <span>高清</span>
                            <span className="text-xs text-gray-400">768×768</span>
                          </div>
                        </Radio>
                        <Radio value="normal" className="text-sm">
                          <div className="flex items-center justify-between w-full">
                            <span>普通</span>
                            <span className="text-xs text-gray-400">512×512</span>
                          </div>
                        </Radio>
                      </Radio.Group>
                    </div>

                    <div>
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-sm font-medium">风格</span>
                        <Tooltip title="不同风格会影响生成效果">
                          <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                        </Tooltip>
                      </div>
                      <Radio.Group 
                        value={style} 
                        onChange={e => setStyle(e.target.value)}
                        className="flex flex-col gap-2"
                      >
                        <Radio value="realistic" className="text-sm">
                          <div className="flex items-center justify-between w-full">
                            <span>写实</span>
                            <span className="text-xs text-gray-400">��实感强</span>
                          </div>
                        </Radio>
                        <Radio value="artistic" className="text-sm">
                          <div className="flex items-center justify-between w-full">
                            <span>艺术</span>
                            <span className="text-xs text-gray-400">创意表现</span>
                          </div>
                        </Radio>
                        <Radio value="anime" className="text-sm">
                          <div className="flex items-center justify-between w-full">
                            <span>动漫</span>
                            <span className="text-xs text-gray-400">二次元风格</span>
                          </div>
                        </Radio>
                      </Radio.Group>
                    </div>
                  </div>

                  {/* 采样步数设置 */}
                  <div>
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">采样步数</span>
                        <Tooltip title="步数越多，细节越丰富，但耗时更长">
                          <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                        </Tooltip>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="flex items-center gap-1">
                          <Button 
                            size="small"
                            type="text"
                            onClick={() => setSteps(Math.max(20, steps - 10))}
                            icon={<MinusOutlined className="text-xs" />}
                          />
                          <Input
                            value={steps}
                            onChange={e => {
                              const value = parseInt(e.target.value);
                              if (value >= 20 && value <= 150) {
                                setSteps(value);
                              }
                            }}
                            className="w-12 text-center text-xs"
                          />
                          <Button
                            size="small"
                            type="text" 
                            onClick={() => setSteps(Math.min(150, steps + 10))}
                            icon={<PlusOutlined className="text-xs" />}
                          />
                        </div>
                        <span className="text-xs text-gray-400">步</span>
                      </div>
                    </div>
                    <Slider
                      min={20}
                      max={150}
                      value={steps}
                      onChange={setSteps}
                      marks={{
                        20: { label: '20', style: { fontSize: '10px', transform: 'translateX(-50%)' } },
                        50: { label: '50', style: { fontSize: '10px', transform: 'translateX(-50%)' } },
                        100: { label: '100', style: { fontSize: '10px', transform: 'translateX(-50%)' } },
                        150: { label: '150', style: { fontSize: '10px', transform: 'translateX(-50%)' } }
                      }}
                    />
                  </div>
                </div>

                {/* 高级参数设置 */}
                {showAdvanced && (
                  <>
                    <Divider className="my-4" />
                    <div className="space-y-6">
                      {/* 提示词相关性设置 */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <span className="text-sm font-medium">提示词相关性 (CFG Scale)</span>
                            <Tooltip title="值越大，生成的图像越接近提示词描述，但可能降低创意性">
                              <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                            </Tooltip>
                          </div>
                          <div className="flex items-center gap-2">
                            <div className="flex items-center gap-1">
                              <Button 
                                size="small"
                                type="text"
                                onClick={() => setCfg(Math.max(1, cfg - 1))}
                                icon={<MinusOutlined className="text-xs" />}
                              />
                              <Input
                                value={cfg}
                                onChange={e => {
                                  const value = parseInt(e.target.value);
                                  if (value >= 1 && value <= 20) {
                                    setCfg(value);
                                  }
                                }}
                                className="w-12 text-center text-xs"
                              />
                              <Button
                                size="small"
                                type="text"
                                onClick={() => setCfg(Math.min(20, cfg + 1))}
                                icon={<PlusOutlined className="text-xs" />}
                              />
                            </div>
                          </div>
                        </div>
                        <Slider
                          min={1}
                          max={20}
                          value={cfg}
                          onChange={setCfg}
                          marks={{
                            1: { label: '创意', style: { fontSize: '10px' } },
                            7: { label: '平衡', style: { fontSize: '10px' } },
                            20: { label: '精确', style: { fontSize: '10px' } }
                          }}
                        />
                      </div>

                      {/* 混乱程度和风格化程度 */}
                      <div className="grid grid-cols-2 gap-6">
                        <div>
                          <div className="flex items-center justify-between mb-2">
                            <div className="flex items-center gap-2">
                              <span className="text-sm font-medium">混乱程度</span>
                              <Tooltip title="影响生成图像的随机性和创意性">
                                <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                              </Tooltip>
                            </div>
                            <span className="text-xs text-gray-400">{chaos}%</span>
                          </div>
                          <Slider
                            min={0}
                            max={100}
                            value={chaos}
                            onChange={setChaos}
                            marks={{
                              0: { label: '稳定', style: { fontSize: '10px' } },
                              100: { label: '混乱', style: { fontSize: '10px' } }
                            }}
                          />
                        </div>

                        <div>
                          <div className="flex items-center justify-between mb-2">
                            <div className="flex items-center gap-2">
                              <span className="text-sm font-medium">风格化程度</span>
                              <Tooltip title="影响艺术风格的强度">
                                <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                              </Tooltip>
                            </div>
                            <span className="text-xs text-gray-400">{stylization}%</span>
                          </div>
                          <Slider
                            min={0}
                            max={100}
                            value={stylization}
                            onChange={setStylization}
                            marks={{
                              0: { label: '写实', style: { fontSize: '10px' } },
                              100: { label: '艺术', style: { fontSize: '10px' } }
                            }}
                          />
                        </div>
                      </div>

                      {/* 种子值设置 */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            <span className="text-sm font-medium">种子值</span>
                            <Tooltip title="相同的种子值会生成相似的图像">
                              <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                            </Tooltip>
                          </div>
                          <Button
                            size="small"
                            type="text"
                            onClick={() => setSeed(-1)}
                            icon={<UndoOutlined className="text-xs" />}
                          >
                            随机
                          </Button>
                        </div>
                        <Input
                          value={seed === -1 ? '' : seed}
                          onChange={e => setSeed(parseInt(e.target.value) || -1)}
                          placeholder="输入种子值或留空随机"
                          className="text-sm"
                          suffix={
                            <Tooltip title="种子值用于复现相同的生成结果">
                              <QuestionCircleOutlined className="text-gray-400 text-xs cursor-help" />
                            </Tooltip>
                          }
                        />
                      </div>
                    </div>
                  </>
                )}
              </div>
            </div>

            {/* 生成按钮区域 */}
            <div className="bg-white rounded-lg shadow-sm p-4 sticky bottom-0 lg:relative">
              <div className="flex gap-2">
                <Button 
                  type="primary" 
                  block
                  size="large"
                  icon={<SendOutlined />}
                  onClick={handleGenerate}
                  loading={isGenerating}
                  className="h-10 shadow-sm hover:shadow-md transition-shadow flex-1"
                >
                  生成图片
                </Button>
                <Tooltip title="保存当前设置">
                  <Button
                    type="default"
                    size="large"
                    icon={<SaveOutlined />}
                    onClick={handleSavePreset}
                    className="h-10"
                  />
                </Tooltip>
              </div>
            </div>
          </div>

          {/* 右侧生成历史区域 */}
          <div className="flex-1 bg-white rounded-lg shadow-sm flex flex-col overflow-hidden min-h-[500px]">
            <div className="px-4 py-3 border-b border-gray-100 flex flex-col lg:flex-row items-start lg:items-center justify-between gap-2">
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  <HistoryOutlined className="text-gray-400" />
                  <span className="text-sm font-medium">生成历史</span>
                  <Tag color="blue" className="text-xs">{generations.length} 张图片</Tag>
                </div>
                <Select
                  defaultValue="all"
                  size="small"
                  style={{ width: 100 }}
                  options={[
                    { label: '全部', value: 'all' },
                    { label: '收藏', value: 'favorite' },
                    { label: '最近', value: 'recent' }
                  ]}
                />
              </div>
              <div className="flex items-center space-x-2 w-full lg:w-auto justify-between lg:justify-end">
                <Radio.Group size="small" value={viewMode} onChange={e => setViewMode(e.target.value)}>
                  <Radio.Button value="grid" className="text-xs">网格</Radio.Button>
                  <Radio.Button value="list" className="text-xs">列表</Radio.Button>
                </Radio.Group>
                <Button 
                  type="text" 
                  size="small"
                  icon={<SaveOutlined className="text-gray-400" />}
                  className="text-sm hover:text-primary-500 hover:bg-primary-50 transition-colors"
                >
                  保存全部
                </Button>
                <Button 
                  type="text" 
                  size="small"
                  danger
                  icon={<DeleteOutlined />}
                  onClick={handleClear}
                  className="text-sm hover:bg-red-50 transition-colors"
                >
                  清空历史
                </Button>
              </div>
            </div>

            <div className="flex-1 p-3 overflow-y-auto">
              <div className={
                viewMode === 'grid' 
                  ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3' 
                  : 'space-y-3'
              }>
                {generations.map((generation) => (
                  <div
                    key={generation.id}
                    className="bg-white rounded-lg shadow-sm overflow-hidden group hover:shadow-md transition-all duration-200 cursor-pointer"
                    onClick={() => {
                      setSelectedImage(generation);
                      setShowFullscreen(true);
                    }}
                  >
                    <div className={viewMode === 'grid' ? 'relative aspect-square' : 'relative h-48'}>
                      {generation.status === 'generating' ? (
                        <div className="absolute inset-0 bg-gray-100 flex items-center justify-center">
                          <LoadingOutlined className="text-xl text-primary-500" />
                        </div>
                      ) : (
                        <div className="relative h-full">
                          <img
                            src={generation.imageUrl}
                            alt={generation.prompt}
                            className="w-full h-full object-cover"
                          />
                          <div className="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors duration-200">
                            <div className="absolute top-2 right-2">
                              <Button
                                type="text"
                                size="small"
                                icon={favorites.includes(generation.id) ? <StarFilled className="text-yellow-500" /> : <StarOutlined className="text-white" />}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  handleToggleFavorite(generation.id);
                                }}
                                className="hover:text-yellow-500"
                              />
                            </div>
                            <div className="absolute inset-x-0 bottom-0 p-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                              <div className="flex items-center justify-between">
                                <div className="flex space-x-1">
                                  <Tooltip title="生成变体">
                                    <Button
                                      type="text"
                                      size="small"
                                      icon={<CopyOutlined className="text-white text-sm" />}
                                      onClick={(e) => {
                                        e.stopPropagation();
                                        handleVariation(generation);
                                      }}
                                      className="hover:text-primary-500"
                                    />
                                  </Tooltip>
                                  <Tooltip title="放大">
                                    <Button
                                      type="text"
                                      size="small"
                                      icon={<PlusOutlined className="text-white text-sm" />}
                                      onClick={(e) => {
                                        e.stopPropagation();
                                        handleUpscale(generation);
                                      }}
                                      className="hover:text-primary-500"
                                    />
                                  </Tooltip>
                                  <Tooltip title="局部重绘">
                                    <Button
                                      type="text"
                                      size="small"
                                      icon={<EditOutlined className="text-white text-sm" />}
                                      onClick={(e) => {
                                        e.stopPropagation();
                                        handleInpaint(generation);
                                      }}
                                      className="hover:text-primary-500"
                                    />
                                  </Tooltip>
                                </div>
                                <Tooltip title="下载">
                                  <Button
                                    type="text"
                                    size="small"
                                    icon={<DownloadOutlined className="text-white text-sm" />}
                                    className="hover:text-primary-500"
                                  />
                                </Tooltip>
                              </div>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                    <div className="p-2 border-t border-gray-100">
                      <p className="text-xs text-gray-700 line-clamp-2 mb-1.5">{generation.prompt}</p>
                      <div className="flex items-center justify-between">
                        <span className="text-xs text-gray-400">
                          {generation.timestamp.toLocaleTimeString()}
                        </span>
                        <Tag color="blue" className="text-xs ml-1">{generation.model}</Tag>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 提示词模板抽屉 */}
      <Drawer
        title="提示词模板"
        placement="right"
        onClose={() => setShowTemplates(false)}
        open={showTemplates}
        width={400}
      >
        {PROMPT_TEMPLATES.map(category => (
          <div key={category.category} className="mb-6">
            <div className="text-sm font-medium mb-2">{category.category}</div>
            <div className="flex flex-wrap gap-2">
              {category.items.map(item => (
                <Tag
                  key={item.label}
                  className="cursor-pointer hover:bg-primary-50"
                  onClick={() => handleApplyTemplate(item.value)}
                >
                  {item.label}
                </Tag>
              ))}
            </div>
          </div>
        ))}
        <Divider />
        <div className="mb-6">
          <div className="text-sm font-medium mb-2">负向提示词模板</div>
          <div className="flex flex-wrap gap-2">
            {NEGATIVE_TEMPLATES.map(item => (
              <Tag
                key={item.label}
                className="cursor-pointer hover:bg-primary-50"
                onClick={() => setNegativePrompt(item.value)}
              >
                {item.label}
              </Tag>
            ))}
          </div>
        </div>
      </Drawer>

      {/* 全屏预览抽屉 */}
      <Drawer
        title={null}
        placement="right"
        onClose={() => {
          setShowFullscreen(false);
          setSelectedImage(null);
        }}
        open={showFullscreen && selectedImage !== null}
        width="100%"
        className="!bg-black/90"
        contentWrapperStyle={{ background: 'transparent', boxShadow: 'none' }}
      >
        <div className="h-full flex items-center justify-center">
          <div className="relative">
            <img
              src={selectedImage?.imageUrl}
              alt={selectedImage?.prompt}
              className="max-h-[90vh] max-w-[90vw] object-contain"
            />
            <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-4">
              <p className="text-white">{selectedImage?.prompt}</p>
              <div className="mt-2 flex items-center justify-between text-white/60">
                <span>{selectedImage?.timestamp.toLocaleString()}</span>
                <span>{selectedImage?.width}x{selectedImage?.height}</span>
              </div>
            </div>
          </div>
        </div>
      </Drawer>
    </div>
  );
};

export default Image; 