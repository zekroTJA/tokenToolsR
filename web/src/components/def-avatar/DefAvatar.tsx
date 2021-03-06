/** @format */

import React, { Component, CSSProperties } from 'react';

import './DefAvatar.scss';

const COLORS = [
  '#e57373',
  '#F06292',
  '#BA68C8',
  '#9575CD',
  '#7986CB',
  '#64B5F6',
  '#4FC3F7',
  '#4DD0E1',
  '#4DB6AC',
  '#81C784',
  '#AED581',
  '#DCE775',
  '#FFD54F',
  '#FFB74D',
  '#FF8A65',
  '#A1887F',
];

export interface Props {
  src: string;
  alt: string;
  text: string;
  width: number;
  height: number;
}

export default class DefAvatar extends Component<Props> {
  public state = {
    isError: false,
  };

  public static defaultProps = {
    alt: '',
    text: null,
  };

  public render() {
    const img = (
      <img
        src={this.props.src}
        width={this.props.width}
        height={this.props.height}
        alt={this.props.alt}
        onError={this.onError.bind(this)}
      />
    );
    const alt = <div className="def-avatar" style={this.style}>{this.character}</div>;

    return this.state.isError ? alt : img;
  }

  private onError() {
    this.setState({ isError: true });
  }

  private get color(): string {
    const c = this.character.charCodeAt(0) - 65;
    return COLORS[c % COLORS.length];
  }

  private get character(): string {
    const text = this.props.text ?? this.props.alt;
    return text.substr(0, 1).toUpperCase();
  }

  private get style(): CSSProperties {
    return {
      backgroundColor: this.color,
      width: this.props.width,
      height: this.props.height,
      lineHeight: this.props.height + 'px',
      fontSize: Math.round((this.props.height + this.props.width) / 4),
    };
  }
}
