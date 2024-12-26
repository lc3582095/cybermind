import React, { useState } from 'react';
import { Input, Button, Radio, InputNumber, Tag, message, Progress } from 'antd';
import {
  DeleteOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  LoadingOutlined,
  CustomerServiceOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons';

const { TextArea } = Input;

interface MusicGeneration {
  id: number;
  title: string;
  lyrics: string;
  timestamp: Date;
  status: 'generating' | 'completed' | 'error';
  progress: number;
  duration: number;
  style: string;
  emotion: string;
}

const Music: React.FC = () => {
  const [title, setTitle] = useState('');
  const [lyrics, setLyrics] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [generations, setGenerations] = useState<MusicGeneration[]>([]);
  const [selectedMode, setSelectedMode] = useState('专业模式');
  const [style, setStyle] = useState('流行');
  const [emotion, setEmotion] = useState('快乐');
  const [duration, setDuration] = useState(30);
  const [seed, setSeed] = useState<string>('');
  const [excludeContent, setExcludeContent] = useState<string[]>([]);
  const [excludeInput, setExcludeInput] = useState('');

  const handleGenerate = () => {
    if (!title.trim() || isGenerating) return;

    const newGeneration: MusicGeneration = {
      id: Date.now(),
      title,
      lyrics,
      timestamp: new Date(),
      status: 'generating',
      progress: 0,
      duration,
      style,
      emotion,
    };

    setGenerations([newGeneration, ...generations]);
    setTitle('');
    setLyrics('');
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
    setExcludeContent(excludeContent.filter((t) => t !== tag));
  };

  return (
    <div className="flex h-screen bg-[#F7F9FB]">
      {/* 左侧导航栏 */}
      <div className="w-[240px] bg-white border-r border-gray-100 flex flex-col h-screen">
        {/* 顶部标题 */}
        <div className="h-14 px-4 flex items-center border-b border-gray-100 flex-shrink-0">
          <span className="text-base font-medium">AI音乐</span>
          <span className="text-xs text-gray-400 ml-2">@musicgen</span>
        </div>

        {/* 中间可滚动区域 */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-4 space-y-4">
            {/* 创作模式选择 */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm">创作模式</span>
                <span className="text-xs text-red-500">*</span>
              </div>
              <div className="space-y-2">
                {['专业模式', '简易模式'].map((mode) => (
                  <div
                    key={mode}
                    onClick={() => setSelectedMode(mode)}
                    className={`w-full px-3 py-2 rounded-lg text-sm cursor-pointer transition-colors ${
                      selectedMode === mode
                        ? 'bg-[#EEF2FF] text-[#6366F1]'
                        : 'bg-white text-gray-500 border border-gray-100'
                    }`}
                  >
                    {mode}
                  </div>
                ))}
              </div>
            </div>

            {/* 歌曲标题 */}
            <div>
              <div className="mb-2">
                <span className="text-sm">歌曲标题</span>
              </div>
              <Input
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="请输入20字以内的歌曲名称"
                className="bg-[#F7F9FB] border-none hover:border-none focus:shadow-none"
                maxLength={20}
              />
            </div>

            {/* 歌曲歌词 */}
            <div>
              <div className="mb-2">
                <span className="text-sm">歌曲歌词</span>
              </div>
              <div className="relative">
                <TextArea
                  value={lyrics}
                  onChange={(e) => setLyrics(e.target.value)}
                  placeholder="请输入歌词内容"
                  autoSize={{ minRows: 3, maxRows: 6 }}
                  className="bg-[#F7F9FB] border-none hover:border-none focus:shadow-none pr-16"
                  maxLength={1200}
                />
                <span className="absolute right-2 bottom-2 text-xs text-gray-400">
                  {lyrics.length} / 1200
                </span>
              </div>
            </div>

            {/* 音乐风格 */}
            <div>
              <div className="mb-2">
                <span className="text-sm">音乐风格</span>
              </div>
              <div className="space-y-2">
                <div className="flex flex-wrap gap-2">
                  {['流行', '电子', '民谣'].map((s) => (
                    <span
                      key={s}
                      onClick={() => setStyle(s)}
                      className={`px-3 py-1 rounded-md text-xs cursor-pointer transition-colors ${
                        style === s
                          ? 'bg-[#EEF2FF] text-[#6366F1]'
                          : 'bg-[#F7F9FB] text-gray-500 hover:bg-gray-100'
                      }`}
                    >
                      {s}
                    </span>
                  ))}
                </div>
                <div className="flex flex-wrap gap-2">
                  {['说唱', '布鲁斯', '古风'].map((s) => (
                    <span
                      key={s}
                      onClick={() => setStyle(s)}
                      className={`px-3 py-1 rounded-md text-xs cursor-pointer transition-colors ${
                        style === s
                          ? 'bg-[#EEF2FF] text-[#6366F1]'
                          : 'bg-[#F7F9FB] text-gray-500 hover:bg-gray-100'
                      }`}
                    >
                      {s}
                    </span>
                  ))}
                </div>
                <div className="flex flex-wrap gap-2">
                  {['R&B', '布鲁斯', '金属'].map((s) => (
                    <span
                      key={s}
                      onClick={() => setStyle(s)}
                      className={`px-3 py-1 rounded-md text-xs cursor-pointer transition-colors ${
                        style === s
                          ? 'bg-[#EEF2FF] text-[#6366F1]'
                          : 'bg-[#F7F9FB] text-gray-500 hover:bg-gray-100'
                      }`}
                    >
                      {s}
                    </span>
                  ))}
                </div>
              </div>
            </div>

            {/* 情感倾向 */}
            <div>
              <div className="mb-2">
                <span className="text-sm">情感倾向</span>
              </div>
              <div className="space-y-2">
                <div className="flex flex-wrap gap-2">
                  {['快乐', '悲伤', '安静'].map((e) => (
                    <span
                      key={e}
                      onClick={() => setEmotion(e)}
                      className={`px-3 py-1 rounded-md text-xs cursor-pointer transition-colors ${
                        emotion === e
                          ? 'bg-[#EEF2FF] text-[#6366F1]'
                          : 'bg-[#F7F9FB] text-gray-500 hover:bg-gray-100'
                      }`}
                    >
                      {e}
                    </span>
                  ))}
                </div>
                <div className="flex flex-wrap gap-2">
                  {['激情', '浪漫', '忧郁'].map((e) => (
                    <span
                      key={e}
                      onClick={() => setEmotion(e)}
                      className={`px-3 py-1 rounded-md text-xs cursor-pointer transition-colors ${
                        emotion === e
                          ? 'bg-[#EEF2FF] text-[#6366F1]'
                          : 'bg-[#F7F9FB] text-gray-500 hover:bg-gray-100'
                      }`}
                    >
                      {e}
                    </span>
                  ))}
                </div>
              </div>
            </div>

            {/* 音乐长度 */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm">音乐长度 (秒)</span>
                <InputNumber
                  value={duration}
                  onChange={(value) => setDuration(value || 30)}
                  min={15}
                  max={300}
                  className="w-20 bg-[#F7F9FB] border-none"
                  size="small"
                  controls={false}
                />
              </div>
            </div>

            {/* 种子 */}
            <div>
              <div className="mb-2">
                <span className="text-sm">种子 (可选)</span>
              </div>
              <Input
                value={seed}
                onChange={(e) => setSeed(e.target.value)}
                placeholder="输入数字作为生成种子"
                className="bg-[#F7F9FB] border-none hover:border-none focus:shadow-none"
              />
            </div>

            {/* 排除内容 */}
            <div>
              <div className="mb-2">
                <span className="text-sm">排除内容</span>
              </div>
              <Input
                value={excludeInput}
                onChange={(e) => setExcludeInput(e.target.value)}
                onPressEnter={handleAddExclude}
                placeholder="输入要排除的元素，回车添加"
                suffix={
                  <Button 
                    type="text" 
                    size="small" 
                    onClick={handleAddExclude}
                    className="text-[#6366F1] hover:text-[#5558E6]"
                  >
                    添加
                  </Button>
                }
                className="bg-[#F7F9FB] border-none hover:border-none focus:shadow-none"
              />
              {excludeContent.length > 0 && (
                <div className="flex flex-wrap gap-1 mt-2">
                  {excludeContent.map((tag) => (
                    <Tag
                      key={tag}
                      closable
                      onClose={() => handleRemoveExclude(tag)}
                      className="bg-[#F7F9FB] border-none text-xs"
                    >
                      {tag}
                    </Tag>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* 底部生成按钮 */}
        <div className="flex-shrink-0 p-4 border-t border-gray-100 bg-white">
          <Button
            type="primary"
            block
            onClick={handleGenerate}
            loading={isGenerating}
            className="bg-[#6366F1] hover:bg-[#5558E6]"
          >
            {isGenerating ? '生成中...' : '生成'}
          </Button>
          <div className="text-xs text-gray-400 text-center mt-2">
            Enter发送 / Shift+Enter换行
          </div>
        </div>
      </div>

      {/* 右侧内容区域 */}
      <div className="flex-1 flex flex-col">
        {/* 顶部标题栏 */}
        <div className="h-14 px-6 bg-white border-b border-gray-100 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium">AI音乐</span>
            <span className="text-xs text-gray-400">@musicgen</span>
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
              <CustomerServiceOutlined className="text-4xl mb-4" />
              <p>暂无音乐创作</p>
            </div>
          ) : (
            <div className="space-y-4">
              {generations.map((generation) => (
                <div
                  key={generation.id}
                  className="bg-white rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow"
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center space-x-3">
                      <Button
                        type="text"
                        shape="circle"
                        icon={<PlayCircleOutlined />}
                        className="text-[#6366F1]"
                      />
                      <div>
                        <h3 className="text-sm font-medium">{generation.title}</h3>
                        <div className="flex items-center space-x-2 text-xs text-gray-400 mt-1">
                          <span>{generation.style}</span>
                          <span>·</span>
                          <span>{generation.emotion}</span>
                          <span>·</span>
                          <span>{generation.duration}s</span>
                        </div>
                      </div>
                    </div>
                    {generation.status === 'generating' ? (
                      <Progress
                        type="circle"
                        percent={generation.progress}
                        width={32}
                        strokeWidth={4}
                        strokeColor={{ '0%': '#818cf8', '100%': '#6366f1' }}
                      />
                    ) : (
                      <Button
                        type="text"
                        icon={<DeleteOutlined />}
                        className="text-gray-400 hover:text-red-500"
                      />
                    )}
                  </div>
                  {generation.lyrics && (
                    <p className="text-sm text-gray-600 bg-gray-50 rounded-lg p-3 mt-2">
                      {generation.lyrics}
                    </p>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Music; 