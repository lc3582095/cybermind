import React, { useState } from 'react';
import { Input, Button, Select, Upload, message, Progress, Slider, Radio, Tag, Tooltip, InputNumber } from 'antd';
import {
  SendOutlined,
  DeleteOutlined,
  UploadOutlined,
  DownloadOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  LoadingOutlined,
  VideoCameraOutlined,
  InfoCircleOutlined,
  RotateLeftOutlined,
  RotateRightOutlined,
  ZoomInOutlined,
  ZoomOutOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import type { UploadFile } from 'antd/es/upload/interface';

const { TextArea } = Input;

interface VideoGeneration {
  id: number;
  prompt: string;
  videoUrl: string;
  thumbnailUrl: string;
  timestamp: Date;
  status: 'generating' | 'completed' | 'error';
  progress: number;
  aspectRatio: string;
  fps: number;
  style: string;
  seed?: number;
}

const Video: React.FC = () => {
  const [prompt, setPrompt] = useState('');
  const [selectedModel, setSelectedModel] = useState('text-to-video');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generations, setGenerations] = useState<VideoGeneration[]>([]);
  const [referenceVideo, setReferenceVideo] = useState<UploadFile | null>(null);
  const [aspectRatio, setAspectRatio] = useState('16:9');
  const [fps, setFps] = useState(30);
  const [style, setStyle] = useState('动漫');
  const [textConsistency, setTextConsistency] = useState(16);
  const [motionEnhancement, setMotionEnhancement] = useState(0);
  const [cameraControl, setCameraControl] = useState({
    horizontal: 0,
    vertical: 0,
    tilt: 0,
  });
  const [excludeContent, setExcludeContent] = useState<string[]>([]);
  const [excludeInput, setExcludeInput] = useState('');
  const [seed, setSeed] = useState<number | null>(null);
  const [zoom, setZoom] = useState(1);
  const [rotation, setRotation] = useState(0);

  const handleGenerate = () => {
    if (!prompt.trim() || isGenerating) return;

    const newGeneration: VideoGeneration = {
      id: Date.now(),
      prompt,
      videoUrl: '',
      thumbnailUrl: 'https://via.placeholder.com/640x360',
      timestamp: new Date(),
      status: 'generating',
      progress: 0,
      aspectRatio,
      fps,
      style,
    };

    setGenerations([newGeneration, ...generations]);
    setPrompt('');
    setIsGenerating(true);

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
                  videoUrl: 'https://example.com/video.mp4', // 模拟视频URL
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

  const handleAddExclude = () => {
    if (excludeInput && !excludeContent.includes(excludeInput)) {
      setExcludeContent([...excludeContent, excludeInput]);
      setExcludeInput('');
    }
  };

  const handleRemoveExclude = (tag: string) => {
    setExcludeContent(excludeContent.filter(t => t !== tag));
  };

  return (
    <div className="flex-1 bg-[#F7F9FB]">
      {/* 顶部标题栏 */}
      <div className="h-14 px-4 bg-white flex items-center justify-between border-b border-gray-100">
        <div className="flex items-center space-x-2">
          <span className="text-base font-medium">AI视频</span>
          <span className="text-xs text-gray-400">@text-to-video</span>
        </div>
        <div className="flex items-center space-x-2">
          <div className="flex items-center space-x-1 mr-4">
            <Button.Group size="small">
              <Button icon={<ZoomOutOutlined />} />
              <Button icon={<ZoomInOutlined />} />
            </Button.Group>
            <Button size="small" icon={<SyncOutlined />}>重置</Button>
          </div>
          <Button type="primary" ghost size="small">作品广场</Button>
        </div>
      </div>

      {/* 主要内容区域 */}
      <div className="flex h-[calc(100vh-3.5rem)]">
        {/* 左侧控制面板 */}
        <div className="w-[280px] bg-white border-r border-gray-100 flex flex-col">
          {/* 设置内容区域 */}
          <div className="flex-1 overflow-y-auto">
            <div className="p-4 space-y-6">
              {/* 上传图片区域 */}
              <div className="border border-dashed border-gray-200 rounded-lg p-4 text-center bg-gray-50">
                <Upload.Dragger
                  accept="image/*"
                  maxCount={1}
                  className="border-none hover:border-none"
                >
                  <p className="text-gray-400">
                    <UploadOutlined className="text-xl mb-1" />
                    <br />
                    点击或拖拽上传参考图片
                  </p>
                </Upload.Dragger>
              </div>

              {/* 视频比例选择 */}
              <div>
                <div className="flex items-center justify-between mb-3">
                  <span className="text-sm">视频比例</span>
                  <Radio.Group value={aspectRatio} onChange={e => setAspectRatio(e.target.value)} size="small">
                    <Radio.Button value="16:9" className="px-2 py-1 text-xs">16:9</Radio.Button>
                    <Radio.Button value="9:16" className="px-2 py-1 text-xs">9:16</Radio.Button>
                    <Radio.Button value="1:1" className="px-2 py-1 text-xs">1:1</Radio.Button>
                    <Radio.Button value="4:5" className="px-2 py-1 text-xs">4:5</Radio.Button>
                  </Radio.Group>
                </div>
              </div>

              {/* 帧率选择 */}
              <div>
                <span className="text-sm block mb-2">帧率 (FPS)</span>
                <Select
                  value={fps}
                  onChange={setFps}
                  options={[
                    { value: 24, label: '24fps - 电影感' },
                    { value: 30, label: '30fps - 标准' },
                    { value: 60, label: '60fps - 流畅' },
                  ]}
                  className="w-full"
                  size="small"
                />
              </div>

              {/* 视频风格 */}
              <div>
                <span className="text-sm block mb-2">视频风格</span>
                <div className="flex flex-wrap gap-2">
                  {['动漫', '电影', '3D', '写实', '自然', '黏土', '酷黑', '炫白'].map(s => (
                    <span
                      key={s}
                      onClick={() => setStyle(s)}
                      className={`px-3 py-1 rounded-full text-xs cursor-pointer transition-colors
                        ${style === s 
                          ? 'bg-[#EEF2FF] text-[#6366F1]' 
                          : 'bg-gray-50 text-gray-500 hover:bg-gray-100'
                        }`}
                    >
                      {s}
                    </span>
                  ))}
                </div>
              </div>

              {/* 文本一致性 */}
              <div>
                <span className="text-sm block mb-2">文本一致性</span>
                <Slider
                  value={textConsistency}
                  onChange={setTextConsistency}
                  min={1}
                  max={25}
                  className="custom-slider"
                  tooltip={{ formatter: value => `${value}` }}
                />
              </div>

              {/* 相机控制 */}
              <div>
                <span className="text-sm block mb-3">相机控制</span>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-xs text-gray-400">左右移动</span>
                      <InputNumber
                        value={cameraControl.horizontal}
                        onChange={v => setCameraControl(prev => ({ ...prev, horizontal: v || 0 }))}
                        min={-100}
                        max={100}
                        size="small"
                        className="w-16"
                        controls={false}
                      />
                    </div>
                    <Slider
                      value={cameraControl.horizontal}
                      onChange={v => setCameraControl(prev => ({ ...prev, horizontal: v }))}
                      min={-100}
                      max={100}
                      className="custom-slider"
                    />
                  </div>
                  <div>
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-xs text-gray-400">前后移动</span>
                      <InputNumber
                        value={cameraControl.vertical}
                        onChange={v => setCameraControl(prev => ({ ...prev, vertical: v || 0 }))}
                        min={-100}
                        max={100}
                        size="small"
                        className="w-16"
                        controls={false}
                      />
                    </div>
                    <Slider
                      value={cameraControl.vertical}
                      onChange={v => setCameraControl(prev => ({ ...prev, vertical: v }))}
                      min={-100}
                      max={100}
                      className="custom-slider"
                    />
                  </div>
                  <div>
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-xs text-gray-400">倾斜角度</span>
                      <InputNumber
                        value={cameraControl.tilt}
                        onChange={v => setCameraControl(prev => ({ ...prev, tilt: v || 0 }))}
                        min={-45}
                        max={45}
                        size="small"
                        className="w-16"
                        controls={false}
                      />
                    </div>
                    <Slider
                      value={cameraControl.tilt}
                      onChange={v => setCameraControl(prev => ({ ...prev, tilt: v }))}
                      min={-45}
                      max={45}
                      className="custom-slider"
                    />
                  </div>
                </div>
              </div>

              {/* 排除内容 */}
              <div>
                <span className="text-sm block mb-2">排除内容</span>
                <Input
                  value={excludeInput}
                  onChange={e => setExcludeInput(e.target.value)}
                  onPressEnter={handleAddExclude}
                  placeholder="输入要排除的元素，回车添加"
                  size="small"
                  suffix={
                    <Button size="small" type="text" onClick={handleAddExclude}>
                      添加
                    </Button>
                  }
                />
                {excludeContent.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-2">
                    {excludeContent.map(tag => (
                      <Tag
                        key={tag}
                        closable
                        onClose={() => handleRemoveExclude(tag)}
                        className="bg-gray-50 border-none text-xs"
                      >
                        {tag}
                      </Tag>
                    ))}
                  </div>
                )}
              </div>

              {/* 种子设置 */}
              <div>
                <span className="text-sm block mb-2">种子设置</span>
                <InputNumber
                  value={seed}
                  onChange={value => setSeed(value)}
                  min={0}
                  max={999999999}
                  className="w-full"
                  size="small"
                  controls={false}
                  placeholder="输入数字作为生成种子"
                />
              </div>
            </div>
          </div>

          {/* 底部生成区域 */}
          <div className="shrink-0 p-4 border-t border-gray-100 bg-white">
            <div className="space-y-3">
              <TextArea
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="描述你想要生成的视频内容..."
                autoSize={{ minRows: 3, maxRows: 6 }}
                className="bg-gray-50 border-none hover:border-none focus:border-none text-sm"
              />
              <div className="flex items-center justify-between">
                <Select
                  value={selectedModel}
                  onChange={setSelectedModel}
                  options={[
                    { value: 'text-to-video', label: '文本生成' },
                    { value: 'image-to-video', label: '图片动画' },
                    { value: 'video-edit', label: '视频编辑' },
                  ]}
                  className="w-28"
                  size="small"
                />
                <Button
                  type="primary"
                  icon={isGenerating ? <LoadingOutlined /> : <SendOutlined />}
                  onClick={handleGenerate}
                  loading={isGenerating}
                  className="bg-[#6366F1] hover:bg-[#5558E6]"
                  size="small"
                >
                  {isGenerating ? '生成中' : '生成'}
                </Button>
              </div>
            </div>
          </div>
        </div>

        {/* 右侧历史创作区域 */}
        <div className="flex-1 bg-white">
          {/* 历史记录标题栏 */}
          <div className="h-12 px-6 border-b border-gray-100 flex items-center justify-between">
            <span className="text-sm font-medium">历史创作</span>
            <div className="flex items-center space-x-2">
              <Button size="small" type="text" icon={<DeleteOutlined />} onClick={handleClear}>
                清空
              </Button>
            </div>
          </div>

          {/* 历史记录列表 */}
          <div className="overflow-y-auto h-[calc(100vh-7.5rem)]">
            {generations.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-400">
                <VideoCameraOutlined className="text-4xl mb-4" />
                <p>开始创作你的第一个视频吧</p>
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-4 p-6">
            {generations.map((generation) => (
              <div
                key={generation.id}
                    className="bg-white rounded-lg shadow-sm overflow-hidden group border border-gray-100 hover:shadow-md transition-shadow"
              >
                <div className="relative aspect-video">
                  {generation.status === 'generating' ? (
                        <div className="absolute inset-0 bg-gray-50 flex flex-col items-center justify-center">
                          <LoadingOutlined className="text-2xl text-[#6366F1] mb-2" />
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
                        src={generation.thumbnailUrl}
                        alt={generation.prompt}
                        className="w-full h-full object-cover"
                      />
                      <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <Button
                          type="text"
                          icon={<PlayCircleOutlined />}
                              className="text-4xl text-white hover:text-[#6366F1]"
                        />
                      </div>
                          <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/60 to-transparent p-3 opacity-0 group-hover:opacity-100 transition-opacity">
                    <div className="flex items-center justify-between">
                      <Button
                        type="text"
                        icon={<DownloadOutlined />}
                                className="text-white hover:text-[#6366F1]"
                                size="small"
                      />
                      <Button
                        type="text"
                        icon={<DeleteOutlined />}
                        className="text-white hover:text-red-500"
                                size="small"
                      />
                    </div>
                  </div>
                        </>
                      )}
                    </div>
                    <div className="p-3">
                      <p className="text-sm text-gray-700 line-clamp-2 mb-2">{generation.prompt}</p>
                      <div className="flex items-center justify-between text-xs text-gray-400">
                        <div className="flex items-center space-x-2">
                          <span>{generation.aspectRatio}</span>
                          <span>·</span>
                          <span>{generation.fps}fps</span>
                          <span>·</span>
                          <span>{generation.style}</span>
                </div>
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
    </div>
  );
};

export default Video; 