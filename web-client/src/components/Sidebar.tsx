import React from 'react';
import { Button, Tooltip } from 'antd';
import { Link, useLocation } from 'react-router-dom';
import {
  MessageOutlined,
  VideoCameraOutlined,
  PictureOutlined,
  FileTextOutlined,
  CustomerServiceOutlined,
  AppstoreOutlined,
  SettingOutlined,
  UserOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';

const NavButton: React.FC<{
  icon: React.ReactNode;
  label: string;
  to?: string;
  active?: boolean;
  onClick?: () => void;
}> = ({ icon, label, to, active, onClick }) => {
  const content = (
    <Button
      type="text"
      icon={icon}
      onClick={onClick}
      className={`w-10 h-10 flex items-center justify-center hover:text-[#6366F1] ${
        active ? 'text-[#6366F1]' : 'text-gray-400'
      }`}
    />
  );

  if (to) {
    return (
      <Tooltip title={label} placement="right">
        <Link to={to}>{content}</Link>
      </Tooltip>
    );
  }

  return <Tooltip title={label} placement="right">{content}</Tooltip>;
};

const Sidebar: React.FC = () => {
  const location = useLocation();
  
  const menuItems = [
    { key: '/chat', icon: <MessageOutlined />, label: '聊天' },
    { key: '/image', icon: <PictureOutlined />, label: '绘画' },
    { key: '/video', icon: <VideoCameraOutlined />, label: '视频' },
    { key: '/music', icon: <CustomerServiceOutlined />, label: '音乐' },
    { key: '/ppt', icon: <FileTextOutlined />, label: 'PPT' },
  ];

  return (
    <div className="w-[60px] bg-white flex flex-col items-center py-4 border-r border-gray-100">
      <div className="w-10 h-10 bg-[#6366F1] rounded-lg flex items-center justify-center mb-8 text-white font-bold">
        П
      </div>
      <div className="flex-1 flex flex-col space-y-8">
        {menuItems.map((item) => (
          <NavButton
            key={item.key}
            icon={item.icon}
            label={item.label}
            to={item.key}
            active={location.pathname === item.key}
          />
        ))}
      </div>
      <div className="mt-auto flex flex-col space-y-4">
        <NavButton icon={<SettingOutlined />} label="设置" />
        <NavButton icon={<UserOutlined />} label="用户" />
      </div>
    </div>
  );
};

export default Sidebar; 