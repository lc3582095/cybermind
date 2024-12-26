import React, { useState } from 'react';
import { Input, Button, Select, Tooltip, Dropdown, message as antMessage, Menu, Avatar } from 'antd';
import {
  SendOutlined,
  DeleteOutlined,
  CopyOutlined,
  SettingOutlined,
  ShareAltOutlined,
  EditOutlined,
  MoreOutlined,
  AudioOutlined,
  PictureOutlined,
  FileOutlined,
  CloseOutlined,
  SearchOutlined,
  SoundOutlined,
  ExportOutlined,
  UnorderedListOutlined,
  DesktopOutlined,
  SyncOutlined,
  RobotOutlined,
  ToolOutlined,
  PlusOutlined,
  HistoryOutlined,
  StarOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons';

const { TextArea } = Input;

interface Message {
  id: number;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  status?: 'sending' | 'sent' | 'error';
  type?: 'text' | 'image' | 'code' | 'audio';
  metadata?: {
    model?: string;
    tokens?: number;
    time?: string;
  };
}

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 1,
      role: 'assistant',
      content: '您好，我是 CyberMind AI 助手，很高兴为您服务！我可以帮助您完��各种任务�����括：\n\n• 回答问题和知识咨询\n• 协助编程和代码分析\n• 文案创作和内容优化\n• 数据分析和处理\n\n请问有什么我可以帮您的吗？',
      timestamp: new Date(),
      type: 'text',
      metadata: {
        model: 'gpt-4',
        tokens: 0,
        time: '0.5s',
      },
    }
  ]);
  const [input, setInput] = useState('');
  const [selectedModel, setSelectedModel] = useState('gpt-4');
  const [isGenerating, setIsGenerating] = useState(false);
  const [chatTitle, setChatTitle] = useState('新对话');
  const [showSuggestions, setShowSuggestions] = useState(true);

  const suggestions = [
    { icon: '📝', title: '文案创作', desc: '帮您撰写专业的文案和内容' },
    { icon: '💻', title: '代码助手', desc: '解决编程问题，优化代码实现' },
    { icon: '🎨', title: '创意设计', desc: '提供设计思路和创意方案' },
    { icon: '📊', title: '数据分析', desc: '协助数据处理和分析工作' },
  ];

  const handleSend = () => {
    if (!input.trim() || isGenerating) return;

    const newMessage: Message = {
      id: Date.now(),
      role: 'user',
      content: input,
      timestamp: new Date(),
      type: 'text',
    };

    setMessages([...messages, newMessage]);
    setInput('');
    setIsGenerating(true);
    setShowSuggestions(false);

    // 模拟API调用
    setTimeout(() => {
      const response: Message = {
        id: Date.now() + 1,
        role: 'assistant',
        content: '这是一个模拟的回复消息。在实际开发中，这里将通过API获取AI的响应。',
        timestamp: new Date(),
        type: 'text',
        metadata: {
          model: selectedModel,
          tokens: 150,
          time: '2.3s',
        },
      };
      setMessages((prev) => [...prev, response]);
      setIsGenerating(false);
    }, 1000);
  };

  return (
    <div className="flex h-screen bg-gray-50">
      {/* 左侧历史记录栏 */}
      <div className="w-64 bg-white border-r border-gray-100 flex flex-col">
        <div className="p-4">
          <Button type="primary" icon={<PlusOutlined />} className="w-full" onClick={() => {
            setMessages([messages[0]]);
            setChatTitle('新对话');
            setShowSuggestions(true);
          }}>
            新对话
          </Button>
        </div>
        <div className="flex-1 overflow-y-auto">
          <div className="px-2">
            <div className="text-xs text-gray-500 px-2 py-1">最近对话</div>
            {/* 历史记录列表 */}
            <div className="space-y-1">
              {[1,2,3].map(i => (
                <div key={i} className="px-2 py-2 hover:bg-gray-100 rounded-lg cursor-pointer">
                  <div className="text-sm font-medium truncate">历史对话 {i}</div>
                  <div className="text-xs text-gray-500 truncate">上次对话时间: 1小时前</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* 主聊天区域 */}
      <div className="flex-1 flex flex-col">
        {/* 顶部标题栏 */}
        <div className="h-14 px-4 border-b border-gray-100 bg-white flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="flex items-center space-x-2">
              <span className="text-base font-medium">{chatTitle}</span>
              <Button
                type="text"
                icon={<EditOutlined />}
                size="small"
                className="opacity-0 group-hover:opacity-100"
                onClick={() => {
                  // TODO: 实现标题编辑功能
                }}
              />
            </div>
            <span className="text-xs px-2 py-1 bg-gray-100 rounded-full text-gray-600">
              {selectedModel.toUpperCase()}
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <Tooltip title="搜索对话">
              <Button type="text" icon={<SearchOutlined />} />
            </Tooltip>
            <Tooltip title="导出对话">
              <Button type="text" icon={<ExportOutlined />} />
            </Tooltip>
            <Tooltip title="分享对话">
              <Button type="text" icon={<ShareAltOutlined />} />
            </Tooltip>
            <Tooltip title="设置">
              <Button type="text" icon={<SettingOutlined />} />
            </Tooltip>
          </div>
        </div>

        {/* 消息区域 */}
        <div className="flex-1 overflow-y-auto">
          <div className="max-w-4xl mx-auto py-4">
            {showSuggestions && (
              <div className="grid grid-cols-2 gap-4 p-4">
                {suggestions.map((item, index) => (
                  <div
                    key={index}
                    className="p-4 bg-white rounded-lg border border-gray-200 hover:border-primary-500 cursor-pointer transition-colors"
                    onClick={() => {
                      setInput(`请帮我${item.desc}`);
                      setShowSuggestions(false);
                    }}
                  >
                    <div className="text-2xl mb-2">{item.icon}</div>
                    <div className="font-medium mb-1">{item.title}</div>
                    <div className="text-sm text-gray-500">{item.desc}</div>
                  </div>
                ))}
              </div>
            )}
            
            {messages.map((message) => (
              <div
                key={message.id}
                className={`py-6 px-4 group transition-colors ${
                  message.role === 'assistant' 
                    ? 'bg-gray-50' 
                    : 'flex flex-row-reverse'
                }`}
              >
                <div className={`max-w-3xl mx-auto flex items-start space-x-4 ${
                  message.role === 'user' ? 'flex-row-reverse space-x-reverse' : ''
                }`}>
                  <Avatar
                    className={message.role === 'assistant' ? 'bg-primary-500' : 'bg-blue-500'}
                    icon={message.role === 'assistant' ? <RobotOutlined /> : 'U'}
                  />
                  <div className={`flex-1 min-w-0 ${
                    message.role === 'user' ? 'items-end' : ''
                  }`}>
                    <div className={`flex items-center justify-between mb-2 ${
                      message.role === 'user' ? 'flex-row-reverse' : ''
                    }`}>
                      <div className="flex items-center space-x-2">
                        <span className="font-medium">
                          {message.role === 'assistant' ? 'CyberMind AI' : '我'}
                        </span>
                        {message.metadata && (
                          <span className="text-xs text-gray-500">
                            {message.metadata.tokens}tokens · {message.metadata.time}
                          </span>
                        )}
                      </div>
                      <div className="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <Button
                          type="text"
                          icon={<CopyOutlined />}
                          size="small"
                          onClick={() => {
                            navigator.clipboard.writeText(message.content);
                            antMessage.success('已复制到剪贴板');
                          }}
                        />
                        {message.role === 'assistant' && (
                          <>
                            <Button
                              type="text"
                              icon={<EditOutlined />}
                              size="small"
                            />
                            <Button
                              type="text"
                              icon={<StarOutlined />}
                              size="small"
                            />
                          </>
                        )}
                      </div>
                    </div>
                    <div className={`text-sm leading-relaxed whitespace-pre-wrap rounded-lg p-4 shadow-sm ${
                      message.role === 'assistant' 
                        ? 'bg-white text-gray-800 border border-gray-100' 
                        : 'bg-primary-500 text-white'
                    }`}>
                      {message.content}
                    </div>
                  </div>
                </div>
              </div>
            ))}
            
            {isGenerating && (
              <div className="py-6 px-4">
                <div className="max-w-3xl mx-auto flex items-start space-x-4">
                  <Avatar className="bg-primary-500" icon={<SyncOutlined spin />} />
                  <div className="flex-1">
                    <div className="text-sm text-gray-600">正在思考...</div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* 底部输入区域 */}
        <div className="border-t border-gray-100 bg-white">
          <div className="max-w-4xl mx-auto px-4 py-3">
            <div className="flex items-center justify-between space-x-2 mb-2">
              <div className="flex items-center space-x-2">
                <Select
                  value={selectedModel}
                  onChange={setSelectedModel}
                  options={[
                    { value: 'gpt-4', label: 'GPT-4' },
                    { value: 'gpt-3.5', label: 'GPT-3.5' },
                    { value: 'claude', label: 'Claude' },
                  ]}
                  className="w-32"
                />
                <Button.Group>
                  <Tooltip title="上传图片">
                    <Button icon={<PictureOutlined />} />
                  </Tooltip>
                  <Tooltip title="上传文件">
                    <Button icon={<FileOutlined />} />
                  </Tooltip>
                  <Tooltip title="语音输入">
                    <Button icon={<AudioOutlined />} />
                  </Tooltip>
                </Button.Group>
              </div>
              
              <div className="flex items-center space-x-2">
                <Tooltip title="快捷指令">
                  <Button icon={<UnorderedListOutlined />} />
                </Tooltip>
                <Tooltip title="清空对话">
                  <Button 
                    danger 
                    icon={<DeleteOutlined />}
                    onClick={() => {
                      setMessages([messages[0]]);
                      antMessage.success('对话已清空');
                    }}
                  />
                </Tooltip>
              </div>
            </div>
            
            <div className="relative">
              <TextArea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="输入消息内容，按 Enter 发送，Shift + Enter 换行"
                autoSize={{ minRows: 1, maxRows: 4 }}
                onPressEnter={(e) => {
                  if (!e.shiftKey) {
                    e.preventDefault();
                    handleSend();
                  }
                }}
                disabled={isGenerating}
                className="pr-24 py-3 rounded-lg resize-none border-gray-200 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 hover:border-primary-500 transition-colors"
              />
              <div className="absolute right-2 bottom-2 flex items-center space-x-2">
                {input.trim().length > 0 && (
                  <Button
                    type="text"
                    icon={<DeleteOutlined />}
                    onClick={() => setInput('')}
                    className="text-gray-400 hover:text-gray-600"
                  />
                )}
                <Button
                  type="primary"
                  icon={isGenerating ? <CloseOutlined /> : <SendOutlined />}
                  onClick={isGenerating ? () => setIsGenerating(false) : handleSend}
                  loading={isGenerating}
                >
                  {isGenerating ? '停止' : '发送'}
                </Button>
              </div>
            </div>
            
            <div className="flex items-center justify-between mt-2 text-xs text-gray-500">
              <div className="flex items-center space-x-4">
                <span>按 Enter 发送，Shift + Enter 换行</span>
                <span>|</span>
                <span>当前输入: {input.length} 字</span>
              </div>
              <div className="flex items-center space-x-4">
                <span>Token 余额: 50000</span>
                <a href="#" className="text-primary-500 hover:text-primary-600 flex items-center space-x-1">
                  <span>购买配额</span>
                  <ExportOutlined className="text-xs" />
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Chat; 